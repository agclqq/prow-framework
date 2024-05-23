package mtd

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/tsdb/fileutil"
)

func GetFileSize(req *http.Request) (int64, error) {
	req.Method = "HEAD"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return GetSize(resp)
}
func GetSize(resp *http.Response) (int64, error) {
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("requesting the file size status code：%d", resp.StatusCode)
	}
	sizeStr := resp.Header.Get("Content-Length")
	return strconv.ParseInt(sizeStr, 10, 64)
}

// FormatBytes 根据字节大小格式化输出
func FormatBytes(bytes int64) string {
	const (
		KB int64 = 1 << 10
		MB int64 = 1 << 20
		GB int64 = 1 << 30
		TB int64 = 1 << 40
	)

	if bytes == 0 {
		return "0B"
	}

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2fTB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(MB))
	default:
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(KB))
	}
}

func CleanFiles(downloadFile string) error {
	dir := filepath.Dir(downloadFile)
	fileName := filepath.Base(downloadFile)
	pattern, err := regexp.Compile(fileName + `(\.part\.[0-9]+\..+|download|download\.cfg)`)
	if err != nil {
		return err
	}
	infos, err := ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if !info.IsDir() && pattern.MatchString(info.Name()) {
			err = os.Remove(dir + string(filepath.Separator) + info.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func PrepareDownloadFile(urlStr, saveDir, saveFileName string) (string, error) {
	downloadFileName := saveFileName
	if downloadFileName == "" {
		downloadFileName = GetFileName(urlStr, "download")
	}
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path.Join(saveDir, downloadFileName), nil
}

// GetFileName 从URL中解析文件名
func GetFileName(urlStr string, dft string) string {
	//url 解析
	u, err := url.Parse(urlStr)
	if err != nil {
		return dft
	}
	// 获取路径最后的部分
	fileName := path.Base(u.Path)
	if fileName == "" || fileName == "/" || fileName == "." {
		return dft
	}
	return fileName
}

func DownloadPart(req *http.Request, start, end int64, chunkFileName string, wg *sync.WaitGroup, chThread chan struct{}) error {
	defer wg.Done()
	defer func() { chThread <- struct{}{} }()
	chunkFile, err := os.Create(chunkFileName)
	if err != nil {
		return err
	}
	defer chunkFile.Close()
	reqN := &http.Request{
		Method:     req.Method,
		URL:        req.URL,
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		Header:     http.Header{},
		Body:       req.Body,
	}
	for k, v := range req.Header {
		reqN.Header.Add(k, v[0])
	}
	// 设置下载范围
	reqN.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
	resp, err := http.DefaultClient.Do(reqN)
	if err != nil {
		errStr := fmt.Sprintf("HTTP请求失败:%v", err)
		fmt.Println(errStr)
		return errors.New(errStr)
	}
	defer resp.Body.Close()

	// 写入文件块
	_, err = io.Copy(chunkFile, resp.Body)
	if err != nil {
		errStr := fmt.Sprintf("写入文件块失败:%v", err)
		fmt.Println(errStr)
		return errors.New(errStr)
	}
	xPart, _, err := ExtractParts(chunkFileName)
	if err != nil {
		return err
	}

	err = fileutil.Rename(chunkFile.Name(), xPart+".done")
	if err != nil {
		return err
	}
	return nil
}

// ExtractParts 提取文件名和文件块编号
func ExtractParts(input string) (string, string, error) {
	// 查找最后一个'.'的位置
	lastDotIndex := strings.LastIndex(input, ".")
	if lastDotIndex == -1 {
		return "", "", errors.New("no '.' found in the input")
	}

	// 提取x.part和y
	xPart := input[:lastDotIndex]
	yPart := input[lastDotIndex+1:]
	return xPart, yPart, nil
}
func MergeFileBlocks(fHandle *os.File, downloadFile string, numBlocks int64) error {
	//如果已完成的文件数不等于块数，说明下载失败
	doneFiles, err := GetDoneFiles(downloadFile)
	if err != nil {
		return err
	}
	if len(doneFiles) != int(numBlocks) {
		return errors.New("下载失败，有部分失败的文件块")
	}
	// 合并文件块
	for i := int64(0); i < int64(len(doneFiles)); i++ {
		chunkFileName := fmt.Sprintf("%s.part.%d.done", downloadFile, i)
		chunkFile, err := os.Open(chunkFileName)
		fmt.Println("合并文件块：", chunkFileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(fHandle, chunkFile)
		if err != nil {
			return err
		}
		chunkFile.Close()
		err = os.Remove(chunkFileName)
		if err != nil {
			return err
		}
	}
	return nil
}
func GetPartFileIndex(partFile string) (int, error) {
	pattern, err := regexp.Compile(`\.part\.([0-9]+)\..+`)
	if err != nil {
		return 0, err
	}
	rs := pattern.FindAllStringSubmatch(partFile, -1)
	if len(rs) == 0 {
		return 0, errors.New("not found")
	}
	return strconv.Atoi(rs[0][1])
}
func GetDoneFiles(originFile string) ([]string, error) {
	pattern, err := regexp.Compile(filepath.Base(originFile) + `\.part\.[0-9]+\.done`)
	if err != nil {
		return nil, err
	}
	return GetFilesByPattern(filepath.Dir(originFile), pattern)
}

func GetFilesByPattern(dir string, reg *regexp.Regexp) ([]string, error) {
	infos, err := ReadDir(dir)
	if err != nil {
		return nil, err
	}
	doneFiles := make([]string, 0)
	for _, info := range infos {
		if !info.IsDir() && reg.MatchString(info.Name()) {
			doneFiles = append(doneFiles, info.Name())
		}
	}
	return doneFiles, nil

}
func ReadDir(dirPath string) ([]os.FileInfo, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	return dir.Readdir(-1)
}

func FileSha256(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 初始化SHA-256哈希计算器
	hasher := sha256.New()

	// 读取文件内容并更新哈希计算器
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// 计算最终的哈希值
	hash := hasher.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hash)
	return hashString, nil
}
