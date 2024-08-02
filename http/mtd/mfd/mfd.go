package mfd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/tsdb/fileutil"

	"github.com/agclqq/prow-framework/http/mtd"
)

var DefaultFileName = "downloaded_file"

// Mf 多临时文件下载
type Mf struct {
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
}
type MfOption func(*Mf)

func NewMf(ops ...MfOption) (*Mf, error) {
	mf := &Mf{wg: &sync.WaitGroup{}}
	for _, op := range ops {
		op(mf)
	}
	if mf.url == "" && mf.req == nil {
		return nil, errors.New("url or req is empty")
	}

	if mf.req == nil {
		req, err := http.NewRequest(mf.method, mf.url, nil)
		if err != nil {
			return nil, err
		}
		if mf.header != nil {
			req.Header = mf.header
		}
		mf.req = req
	}
	if mf.req.Header == nil {
		mf.req.Header = make(http.Header)
	}

	if mf.concurrence <= 0 {
		mf.concurrence = 1
	}
	if mf.blockSize <= 0 {
		mf.blockSize = int64(syscall.Getpagesize()) * 1024
	}
	return mf, nil
}

func WithUrl(url string) MfOption {
	return func(m *Mf) {
		m.url = url
	}
}
func WithReq(req *http.Request) MfOption {
	return func(m *Mf) {
		m.req = req
	}
}
func WithSavePath(savePath string) MfOption {
	return func(m *Mf) {
		m.savePath = savePath
	}
}
func WithSaveFileName(saveFileName string) MfOption {
	return func(m *Mf) {
		m.saveFileName = saveFileName
	}
}
func WithConcurrence(threadNum int) MfOption {
	return func(m *Mf) {
		m.concurrence = threadNum
	}
}
func WithBlockSize(blockSize int64) MfOption {
	return func(m *Mf) {
		m.blockSize = blockSize
	}
}

func (mf *Mf) Download() (string, error) {
	return mf.download(false)
}
func (mf *Mf) Resume() (string, error) {
	return mf.download(true)
}
func (mf *Mf) download(resume bool) (string, error) {
	startTime := time.Now()
	urlStr := mf.req.URL.String()
	fmt.Println(filepath.Abs(mf.savePath))
	fileSize, err := mf.getFileSize()
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
	numBlocks := (fileSize + mf.blockSize - 1) / mf.blockSize
	if int64(mf.concurrence) >= numBlocks {
		mf.concurrence = int(numBlocks)
	}

	// 准备下载文件
	downloadFile, err := mtd.PrepareDownloadFile(urlStr, mf.savePath, mf.saveFileName)
	if err != nil {
		return "", err
	}

	//并发控制
	mf.chThread = make(chan struct{}, mf.concurrence)
	defer close(mf.chThread)
	for i := 0; i < mf.concurrence; i++ {
		mf.chThread <- struct{}{}
	}

	doneFilesMap := make(map[int]struct{})
	if resume {
		doneFiles, err := mtd.GetDoneFiles(downloadFile)
		if err != nil {
			return "", err
		}
		for _, doneFile := range doneFiles {
			index, err := mtd.GetPartFileIndex(doneFile)
			if err != nil {
				return "", err
			}
			doneFilesMap[index] = struct{}{}
		}
	} else {
		// 清理已存在的文件
		err = mtd.CleanFiles(downloadFile)
		if err != nil {
			return "", err
		}
	}

	// 启动多个协程下载文件块
	for i := int64(0); i < numBlocks; i++ {
		if _, ok := doneFilesMap[int(i)]; ok && resume {
			continue
		}
		mf.wg.Add(1)
		start := i * mf.blockSize
		end := (i + 1) * mf.blockSize
		if end > fileSize {
			end = fileSize
		}
		chunkFileName := fmt.Sprintf("%s.part.%d.todo", downloadFile, i)
		if err != nil {
			return "", err
		}
		_ = <-mf.chThread
		go func() {
			err = mtd.DownloadPart(mf.req, start, end, chunkFileName, mf.wg, mf.chThread)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	// 等待所有协程完成
	mf.wg.Wait()
	// 合并文件块
	// 创建临时文件
	fHandle, err := os.Create(downloadFile + ".download")
	if err != nil {
		errStr := fmt.Sprintf("创建文件失败:%v", err)
		fmt.Println(errStr)
		return "", errors.New(errStr)
	}
	defer fHandle.Close()
	err = mtd.MergeFileBlocks(fHandle, downloadFile, numBlocks)
	if err != nil {
		os.Remove(fHandle.Name())
		return "", err
	}
	err = fileutil.Rename(fHandle.Name(), downloadFile)
	return downloadFile, err
}

func (mf *Mf) getFileSize() (int64, error) {
	oldMethod := mf.req.Method
	defer func() {
		mf.req.Method = oldMethod
	}()
	return mtd.GetFileSize(mf.req)
}

//
//// formatBytes 根据字节大小格式化输出
//func formatBytes(bytes int64) string {
//	const (
//		KB int64 = 1 << 10
//		MB int64 = 1 << 20
//		GB int64 = 1 << 30
//		TB int64 = 1 << 40
//	)
//
//	if bytes == 0 {
//		return "0B"
//	}
//
//	switch {
//	case bytes >= TB:
//		return fmt.Sprintf("%.2fTB", float64(bytes)/float64(TB))
//	case bytes >= GB:
//		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(GB))
//	case bytes >= MB:
//		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(MB))
//	default:
//		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(KB))
//	}
//}
//
//func cleanFiles(downloadFile string) error {
//	dir := filepath.Dir(downloadFile)
//	fileName := filepath.Base(downloadFile)
//	pattern, err := regexp.Compile(fileName + `(\.part\.[0-9]+\..+|download)`)
//	if err != nil {
//		return err
//	}
//	infos, err := readDir(dir)
//	if err != nil {
//		return err
//	}
//	for _, info := range infos {
//		if !info.IsDir() && pattern.MatchString(info.Name()) {
//			err = os.Remove(dir + string(filepath.Separator) + info.Name())
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//func prepareDownloadFile(urlStr, saveDir, saveFileName string) (string, error) {
//	downloadFileName := saveFileName
//	if downloadFileName == "" {
//		downloadFileName = getFileName(urlStr)
//	}
//	err := os.MkdirAll(saveDir, os.ModePerm)
//	if err != nil {
//		return "", err
//	}
//	return path.Join(saveDir, downloadFileName), nil
//}
//
//// 从URL中解析文件名
//func getFileName(urlStr string) string {
//	//url 解析
//	u, err := url.Parse(urlStr)
//	if err != nil {
//		return DefaultFileName
//	}
//	// 获取路径最后的部分
//	fileName := path.Base(u.Path)
//	if fileName == "" || fileName == "/" || fileName == "." {
//		return DefaultFileName
//	}
//	return fileName
//}
//
//func getSize(resp *http.Response) (int64, error) {
//	if resp.StatusCode != http.StatusOK {
//		return 0, fmt.Errorf("requesting the file size status code：%d", resp.StatusCode)
//	}
//	sizeStr := resp.Header.GetInfo("Content-Length")
//	return strconv.ParseInt(sizeStr, 10, 64)
//}
//
//func downloadPart(req *http.Request, start, end int64, chunkFileName string, wg *sync.WaitGroup, chThread chan struct{}) error {
//	defer wg.Done()
//	defer func() { chThread <- struct{}{} }()
//	chunkFile, err := os.Create(chunkFileName)
//	if err != nil {
//		return err
//	}
//	defer chunkFile.Close()
//	reqN := &http.Request{
//		Method:     req.Method,
//		URL:        req.URL,
//		Proto:      req.Proto,
//		ProtoMajor: req.ProtoMajor,
//		ProtoMinor: req.ProtoMinor,
//		Header:     http.Header{},
//		Body:       req.Body,
//	}
//	for k, v := range req.Header {
//		reqN.Header.Add(k, v[0])
//	}
//	// 设置下载范围
//	reqN.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
//	resp, err := http.DefaultClient.Do(reqN)
//	if err != nil {
//		errStr := fmt.Sprintf("HTTP请求失败:%v", err)
//		fmt.Println(errStr)
//		return errors.New(errStr)
//	}
//	defer resp.Body.Close()
//
//	// 写入文件块
//	_, err = io.Copy(chunkFile, resp.Body)
//	if err != nil {
//		errStr := fmt.Sprintf("写入文件块失败:%v", err)
//		fmt.Println(errStr)
//		return errors.New(errStr)
//	}
//	xPart, _, err := extractParts(chunkFileName)
//	if err != nil {
//		return err
//	}
//
//	err = fileutil.Rename(chunkFile.Name(), xPart+".done")
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// 提取文件名和文件块编号
//func extractParts(input string) (string, string, error) {
//	// 查找最后一个'.'的位置
//	lastDotIndex := strings.LastIndex(input, ".")
//	if lastDotIndex == -1 {
//		return "", "", errors.New("no '.' found in the input")
//	}
//
//	// 提取x.part和y
//	xPart := input[:lastDotIndex]
//	yPart := input[lastDotIndex+1:]
//	return xPart, yPart, nil
//}
//func mergeFileBlocks(fHandle *os.File, downloadFile string, numBlocks int64) error {
//	//如果已完成的文件数不等于块数，说明下载失败
//	doneFiles, err := getDoneFiles(downloadFile)
//	if err != nil {
//		return err
//	}
//	if len(doneFiles) != int(numBlocks) {
//		return errors.New("下载失败，有部分失败的文件块")
//	}
//	// 合并文件块
//	for i := int64(0); i < int64(len(doneFiles)); i++ {
//		chunkFileName := fmt.Sprintf("%s.part.%d.done", downloadFile, i)
//		chunkFile, err := os.Open(chunkFileName)
//		fmt.Println("合并文件块：", chunkFileName)
//		if err != nil {
//			return err
//		}
//		_, err = io.Copy(fHandle, chunkFile)
//		if err != nil {
//			return err
//		}
//		chunkFile.Close()
//		err = os.Remove(chunkFileName)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//func getPartFileIndex(partFile string) (int, error) {
//	pattern, err := regexp.Compile(`\.part\.([0-9]+)\..+`)
//	if err != nil {
//		return 0, err
//	}
//	rs := pattern.FindAllStringSubmatch(partFile, -1)
//	if len(rs) == 0 {
//		return 0, errors.New("not found")
//	}
//	return strconv.Atoi(rs[0][1])
//}
//func getDoneFiles(originFile string) ([]string, error) {
//	pattern, err := regexp.Compile(filepath.Base(originFile) + `\.part\.[0-9]+\.done`)
//	if err != nil {
//		return nil, err
//	}
//	return getFilesByPattern(filepath.Dir(originFile), pattern)
//}
//
//func getFilesByPattern(dir string, reg *regexp.Regexp) ([]string, error) {
//	infos, err := readDir(dir)
//	if err != nil {
//		return nil, err
//	}
//	doneFiles := make([]string, 0)
//	for _, info := range infos {
//		if !info.IsDir() && reg.MatchString(info.Name()) {
//			doneFiles = append(doneFiles, info.Name())
//		}
//	}
//	return doneFiles, nil
//
//}
//func readDir(dirPath string) ([]os.FileInfo, error) {
//	dir, err := os.Open(dirPath)
//	if err != nil {
//		return nil, err
//	}
//	defer dir.Close()
//	return dir.Readdir(-1)
//}
