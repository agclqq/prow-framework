package stream

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// downloadFromArtifactory 函数从Artifactory下载文件，并返回一个Reader接口
func downloadFromArtifactory(artifactoryURL, filePath string) (io.ReadCloser, error) {
	// 创建Artifactory的URL
	u, err := url.Parse(artifactoryURL)
	if err != nil {
		return nil, err
	}
	u.Path = filePath
	artifactoryURL = u.String()

	// 发起HTTP GET请求到Artifactory获取文件
	resp, err := http.Get(artifactoryURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("artifactory returned non-200 status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// downloadFileHandler 函数处理下载请求，使用流式传输
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 定义Artifactory的URL和文件路径
	artifactoryURL := "http://your-internal-artifactory-url"
	filePath := "/path/to/your/file/a"

	// 从Artifactory下载文件
	reader, err := downloadFromArtifactory(artifactoryURL, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	// 设置响应头信息
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"a\"")

	// 使用io.Copy将文件内容从Artifactory直接复制到HTTP响应中
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//-------------------------

// uploadToArtifactory 函数负责将文件上传到Artifactory
func uploadToArtifactory(artifactoryURL, filePath, targetPath string, file io.Reader) error {
	// 创建multipart表单，用于文件上传
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	// 发送POST请求到Artifactory
	req, err := http.NewRequest("POST", artifactoryURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 如果Artifactory需要认证，可以在这里设置认证信息
	// req.SetBasicAuth("username", "password")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("artifactory returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}

// uploadFileHandler 函数处理文件上传请求
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 从请求中获取文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from request", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 定义Artifactory的URL和目标路径
	artifactoryURL := "http://your-internal-artifactory-url/api/storage/your-repo"
	targetPath := "/path/in/artifactory/to/upload/a"

	// 上传文件到Artifactory
	err = uploadToArtifactory(artifactoryURL, header.Filename, targetPath, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 文件上传成功，返回响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
