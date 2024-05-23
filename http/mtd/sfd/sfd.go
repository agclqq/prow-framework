package sfd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/tsdb/fileutil"

	"github.com/agclqq/prow-framework/http/mtd"
)

// Sf 单临时文件下载
type Sf struct {
	url          string
	method       string
	header       http.Header
	req          *http.Request
	savePath     string
	saveFileName string
	concurrence  int
	blockSize    int64
	chThread     chan struct{}
	wg           *sync.WaitGroup
	lock         *sync.Mutex
}

const cfgPartStartLine = 4

type SfOption func(*Sf)

func NewSfd(opts ...SfOption) (*Sf, error) {
	sf := &Sf{wg: &sync.WaitGroup{}, lock: &sync.Mutex{}}
	for _, opt := range opts {
		opt(sf)
	}
	if sf.url == "" && sf.req == nil {
		return nil, errors.New("url or req is empty")
	}

	if sf.req == nil {
		req, err := http.NewRequest(sf.method, sf.url, nil)
		if err != nil {
			return nil, err
		}
		if sf.header != nil {
			req.Header = sf.header
		}
		sf.req = req
	}
	if sf.req.Header == nil {
		sf.req.Header = make(http.Header)
	}

	if sf.concurrence <= 0 {
		sf.concurrence = 1
	}
	if sf.blockSize <= 0 {
		sf.blockSize = int64(syscall.Getpagesize()) * 1024
	}
	return sf, nil
}

func WithUrl(url string) SfOption {
	return func(s *Sf) {
		s.url = url
	}
}

func WithReq(req *http.Request) SfOption {
	return func(s *Sf) {
		s.req = req
	}
}

func WithSavePath(savePath string) SfOption {
	return func(s *Sf) {
		s.savePath = savePath
	}
}

func WithSaveFileName(saveFileName string) SfOption {
	return func(s *Sf) {
		s.saveFileName = saveFileName
	}
}

func WithConcurrence(concurrence int) SfOption {
	return func(s *Sf) {
		s.concurrence = concurrence
	}
}

func WithBlockSize(blockSize int64) SfOption {
	return func(s *Sf) {
		s.blockSize = blockSize
	}
}
func (sf *Sf) Download() (string, error) {
	return sf.download(false)
}
func (sf *Sf) Resume() (string, error) {
	return sf.download(true)
}

func (sf *Sf) download(resume bool) (string, error) {
	startTime := time.Now()
	urlStr := sf.req.URL.String()
	fileSize, err := sf.getFileSize()
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to get the file size, trying to degrade to single-thread download")
		fileSize = 1
	} else {
		fmt.Println("文件大小：", mtd.FormatBytes(fileSize))
	}
	defer func() {
		costTime := time.Now().Sub(startTime).Seconds()
		fmt.Printf("download completed, time used: %.2f seconds, average speed: %s/s\n", costTime, mtd.FormatBytes(int64(float64(fileSize)/costTime)))
	}()

	// 计算需要多少个块
	numBlocks := (fileSize + sf.blockSize - 1) / sf.blockSize
	if int64(sf.concurrence) >= numBlocks {
		sf.concurrence = int(numBlocks)
	}

	// 准备下载文件
	downloadFile, err := mtd.PrepareDownloadFile(urlStr, sf.savePath, sf.saveFileName)
	if err != nil {
		return "", err
	}
	tmpDownloadFile := downloadFile + ".download"
	tmpCfgFile := downloadFile + ".cfg"

	//并发控制
	sf.chThread = make(chan struct{}, sf.concurrence)
	defer close(sf.chThread)
	for i := 0; i < sf.concurrence; i++ {
		sf.chThread <- struct{}{}
	}

	doneFilesMap := make(map[int]struct{})
	var fHandle *os.File
	var cfgHandle *os.File
	if resume {
		fHandle, err = os.OpenFile(tmpDownloadFile, os.O_RDWR, 0666)
		if err != nil {
			return "", err
		}
		cfgHandle, err = os.OpenFile(tmpCfgFile, os.O_RDWR, 0666)
		if err != nil {
			return "", err
		}
		doneFilesMap, err = sf.getDoneChunks(cfgHandle)
		if err != nil {
			return "", err
		}
	} else {
		// 清理已存在的文件
		err = mtd.CleanFiles(downloadFile)
		if err != nil {
			return "", err
		}
		// 创建下载文件
		fHandle, err = os.Create(tmpDownloadFile)
		if err != nil {
			errStr := fmt.Sprintf("创建文件失败:%v", err)
			fmt.Println(errStr)
			return "", errors.New(errStr)
		}
		defer fHandle.Close()
		// 填充0
		err = sf.fillZero(fHandle, fileSize)
		if err != nil {
			return "", err
		}
		// 创建下载配置文件
		cfgHandle, err = os.Create(tmpCfgFile)
		if err != nil {
			return "", err
		}
		defer cfgHandle.Close()
		err = sf.createCfgFile(cfgHandle, sf.blockSize, numBlocks, sf.concurrence)
	}

	type downloadPart struct {
		index   int
		content []byte
	}

	partErrs := errors.Join()
	// 启动多个协程下载文件块
	for i := int64(0); i < numBlocks; i++ {
		if _, ok := doneFilesMap[int(i)]; ok && resume {
			continue
		}
		sf.wg.Add(1)
		start := i * sf.blockSize
		end := (i + 1) * sf.blockSize
		if end > fileSize {
			end = fileSize
		}

		_ = <-sf.chThread
		go func() {
			err = sf.downloadPart(sf.req, i, start, end, fHandle, cfgHandle, sf.wg, sf.chThread)
			if err != nil {
				partErrs = errors.Join(partErrs, err)
			}
		}()
	}

	// 等待所有协程完成
	sf.wg.Wait()

	if len(partErrs.Error()) > 0 {
		return "", partErrs
	}
	// 合并文件块
	// 创建临时文件

	err = fileutil.Rename(fHandle.Name(), downloadFile)
	return downloadFile, err
}

func (sf *Sf) getFileSize() (int64, error) {
	oldMethod := sf.req.Method
	defer func() {
		sf.req.Method = oldMethod
	}()
	return mtd.GetFileSize(sf.req)
}

func (sf *Sf) createCfgFile(cfgFile *os.File, blockSize, numBlocks int64, concurrence int) error {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d\n", blockSize))
	sb.WriteString(fmt.Sprintf("%d\n", numBlocks))
	sb.WriteString(fmt.Sprintf("%d\n", concurrence))
	for i := int64(0); i < numBlocks; i++ {
		sb.WriteString(fmt.Sprintf("%d 0\n", i))
	}
	_, err := cfgFile.WriteString(sb.String())
	return err
}

func (sf *Sf) fillZero(f *os.File, size int64) error {
	// 创建一个零字节的Reader，并复制fileSize字节到文件
	zeroes := io.NopCloser(io.LimitReader(os.Stdin, size))
	// 将零字节写入文件，直到达到fileSize大小
	_, err := io.Copy(f, zeroes)
	if err != nil {
		return err
	}
	// 确保文件大小正确
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() != size {
		return errors.New("file size does not match the expected size")
	}
	return nil
}

func (sf *Sf) getDoneChunks(cfg *os.File) (map[int]struct{}, error) {
	rs := make(map[int]struct{})
	_, err := cfg.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(cfg)
	for i := 1; scanner.Scan(); i++ {
		if i >= cfgPartStartLine {
			s1, s2, found := strings.Cut(scanner.Text(), " ")
			if !found {
				continue
			}
			index, err := strconv.Atoi(s1)
			if err != nil {
				return nil, err
			}
			status, err := strconv.Atoi(s2)
			if err != nil {
				return nil, err
			}
			if status == 1 {
				rs[index] = struct{}{}
			}
		}
	}
	return rs, nil
}

func (sf *Sf) downloadPart(req *http.Request, lineNum, start, end int64, fHandle, cfgHandle *os.File, wg *sync.WaitGroup, chThread chan struct{}) error {
	defer func() {
		wg.Done()
		chThread <- struct{}{}
	}()

	// 打开文件
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sf.lock.Lock()
	defer sf.lock.Unlock()
	// 写入文件块
	_, err = cfgHandle.Seek(start, 0)
	if err != nil {
		return err
	}
	_, err = io.CopyN(fHandle, resp.Body, end-start)
	if err != nil {
		return err
	}

	//更新配置文件
	scanner := bufio.NewScanner(cfgHandle)
	foundLine := false
	for i := int64(1); scanner.Scan(); i++ {
		if i == lineNum+cfgPartStartLine {
			foundLine = true
			s1, _, found := strings.Cut(scanner.Text(), " ")
			if !found {
				_, err = cfgHandle.WriteString(fmt.Sprintf("%d 1\n", lineNum))
				if err != nil {
					return err
				}
				return nil
			}
			if s1 != fmt.Sprintf("%d", lineNum) {
				return errors.New("line number does not match，need to re-download")
			}
			_, err = cfgHandle.WriteString(fmt.Sprintf("%d 1\n", lineNum))
			if err != nil {
				return err
			}
		}
	}
	if !foundLine {
		return errors.New("line number not found")
	}
	return nil
}
