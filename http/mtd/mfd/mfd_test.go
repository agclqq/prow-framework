package mfd

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/agclqq/prow-framework/http/mtd"
)

type testMf struct {
	Mf
	sha256 string
}

var tmpFilePath = "/tmp/mfd_test/file_path"
var tmpSavePath = "/tmp/mfd_test/save_path"
var files = []string{"url.log.log", "req.log.log", "header.log.log", "body.log.log", "resume.url.log.log"}

func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/file/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		filePath := filepath.Join(tmpFilePath, id)
		//file, err := os.Open(filePath)
		//if err != nil {
		//	w.WriteHeader(http.StatusNotFound)
		//	return
		//}
		//defer file.Close()
		//
		//// 获取文件信息
		//fileInfo, err := file.Stat()
		//if err != nil {
		//	http.Error(w, "无法获取文件信息", http.StatusInternalServerError)
		//	return
		//}

		//w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
		//w.Header().Set("Content-Type", "application/octet-stream")
		//w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		// 将文件内容写入响应
		http.ServeFile(w, r, filePath)
	})
	return mux
}
func svr() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router(),
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sign := <-ch
	fmt.Println("got a sign:", sign)
}
func createTmpFile(fileName string) (string, error) {
	err := os.MkdirAll(tmpFilePath, os.ModePerm)
	if err != nil {
		fmt.Println("create dir err:", err)
		return "", err
	}
	filePath := filepath.Join(tmpFilePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("create file err:", err)
		return "", err
	}
	defer file.Close()
	min := 1024
	max := 5242880
	for i := 0; i < 5; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := r.Intn(max-min+1) + min
		sb := strings.Builder{}
		for j := 0; j < randomNumber; j++ {
			sb.WriteString("a")
		}
		file.WriteString(sb.String())
	}
	return filePath, nil
}
func removeTmpFiles() {
	os.RemoveAll(tmpFilePath)
	os.RemoveAll(tmpSavePath)
}
func createTmpFiles(files []string) ([]string, error) {
	sha256s := make([]string, 0)
	for _, file := range files {
		tmpFile, err := createTmpFile(file)
		if err != nil {
			return nil, err
		}
		sha256, err := mtd.FileSha256(tmpFile)
		if err != nil {
			return nil, err
		}
		sha256s = append(sha256s, sha256)
	}
	return sha256s, nil
}
func TestMf_Download(t *testing.T) {
	t.Parallel()
	go func() {
		svr()
	}()
	time.Sleep(1 * time.Second)

	sha256s, err := createTmpFiles(files)
	if err != nil {
		t.Error(err)
		return
	}
	defer removeTmpFiles()
	var tests = []struct {
		name    string
		mf      testMf
		wantErr bool
	}{
		{name: "t1_dynamicText_url",
			mf: testMf{
				Mf: Mf{
					url:          "http://127.0.0.1:8080/file/" + files[0],
					method:       "",
					header:       nil,
					req:          nil,
					savePath:     tmpSavePath,
					saveFileName: files[0],
					concurrence:  5,
					blockSize:    0,
				},
				sha256: sha256s[0],
			},
			wantErr: false},
		{name: "t1_dynamicText_req", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[0],
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     tmpSavePath,
				saveFileName: files[0],
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "",
		}, wantErr: false},
		{name: "t2_bin_url", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[1],
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     tmpSavePath,
				saveFileName: files[1],
				concurrence:  -1,
				blockSize:    0,
			},
			sha256: sha256s[1],
		}, wantErr: false},
		{name: "t2_bin_req", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[2],
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     tmpSavePath,
				saveFileName: files[2],
				concurrence:  5,
				blockSize:    0,
			},
			sha256: sha256s[2],
		}, wantErr: false},
		{name: "t3_img_url", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[3],
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     tmpSavePath,
				saveFileName: "",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: sha256s[3],
		}, wantErr: false},
		{name: "t3_img_req", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[4],
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     tmpSavePath,
				saveFileName: files[4],
				concurrence:  5,
				blockSize:    0,
			},
			sha256: sha256s[4],
		}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := make([]MfOption, 0)

			if tt.mf.Mf.url != "" {
				opts = append(opts, WithUrl(tt.mf.Mf.url))
			}
			if tt.mf.Mf.req != nil {
				u, err := url.Parse(tt.mf.Mf.url)
				if err != nil {
					t.Error(err)
					return
				}
				tt.mf.Mf.req.URL = u
				tt.mf.Mf.req.Method = tt.mf.Mf.method
				if tt.mf.Mf.req.Method == "" {
					tt.mf.Mf.req.Method = "GET"
				}
				opts = append(opts, WithReq(tt.mf.Mf.req))
			}
			if tt.mf.Mf.savePath != "" {
				opts = append(opts, WithSavePath(tt.mf.Mf.savePath))
			}
			if tt.mf.Mf.saveFileName != "" {
				opts = append(opts, WithSaveFileName(tt.mf.Mf.saveFileName))
			}
			if tt.mf.Mf.concurrence > 0 {
				opts = append(opts, WithConcurrence(tt.mf.Mf.concurrence))
			}
			if tt.mf.Mf.blockSize > 0 {
				opts = append(opts, WithBlockSize(tt.mf.Mf.blockSize))
			}
			mf, err := NewMf(opts...)
			if err != nil {
				t.Error(err)
				return
			}
			downloadFile := ""
			if downloadFile, err = mf.Download(); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.mf.sha256 != "" {
				if err != nil {
					return
				}
				sha256Str, err := mtd.FileSha256(downloadFile)
				if err != nil {
					t.Error(err)
					return
				}
				if sha256Str != tt.mf.sha256 {
					t.Errorf("sha256Str = %v, want %v", sha256Str, tt.mf.sha256)
				}
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			opts := make([]MfOption, 0)

			if tt.mf.Mf.url != "" {
				opts = append(opts, WithUrl(tt.mf.Mf.url))
			}
			if tt.mf.Mf.req != nil {
				u, err := url.Parse(tt.mf.Mf.url)
				if err != nil {
					t.Error(err)
					return
				}
				tt.mf.Mf.req.URL = u
				tt.mf.Mf.req.Method = tt.mf.Mf.method
				if tt.mf.Mf.req.Method == "" {
					tt.mf.Mf.req.Method = "GET"
				}
				opts = append(opts, WithReq(tt.mf.Mf.req))
			}
			if tt.mf.Mf.savePath != "" {
				opts = append(opts, WithSavePath(tt.mf.Mf.savePath))
			}
			if tt.mf.Mf.saveFileName != "" {
				opts = append(opts, WithSaveFileName(tt.mf.Mf.saveFileName))
			}
			if tt.mf.Mf.concurrence > 0 {
				opts = append(opts, WithConcurrence(tt.mf.Mf.concurrence))
			}
			if tt.mf.Mf.blockSize > 0 {
				opts = append(opts, WithBlockSize(tt.mf.Mf.blockSize))
			}
			mf, err := NewMf(opts...)
			if err != nil {
				t.Error(err)
				return
			}
			downloadFile := ""
			if downloadFile, err = mf.Resume(); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.mf.sha256 != "" {
				if err != nil {
					return
				}
				sha256Str, err := mtd.FileSha256(downloadFile)
				if err != nil {
					t.Error(err)
					return
				}
				if sha256Str != tt.mf.sha256 {
					t.Errorf("sha256Str = %v, want %v", sha256Str, tt.mf.sha256)
				}
			}
		})
	}
}

func TestMf_Resume(t *testing.T) {
	t.Parallel()
	go func() {
		svr()
	}()
	time.Sleep(1 * time.Second)

	sha256s, err := createTmpFiles(files)
	if err != nil {
		t.Error(err)
		return
	}
	defer removeTmpFiles()
	var tests = []struct {
		name    string
		mf      testMf
		wantErr bool
	}{
		{name: "t1",
			mf: testMf{
				Mf: Mf{
					url:          "http://127.0.0.1:8080/file/" + files[0],
					method:       "",
					header:       nil,
					req:          nil,
					savePath:     tmpSavePath,
					saveFileName: files[0],
					concurrence:  5,
					blockSize:    0,
				},
				sha256: sha256s[0],
			},
			wantErr: false},
		{name: "t2", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[1],
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     tmpSavePath,
				saveFileName: files[1],
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "",
		}, wantErr: false},
		{name: "t3", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[2],
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     tmpSavePath,
				saveFileName: files[2],
				concurrence:  -1,
				blockSize:    0,
			},
			sha256: sha256s[2],
		}, wantErr: false},
		{name: "t4", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[3],
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     tmpSavePath,
				saveFileName: files[3],
				concurrence:  5,
				blockSize:    0,
			},
			sha256: sha256s[3],
		}, wantErr: false},
		{name: "t5", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/file/" + files[4],
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     tmpSavePath,
				saveFileName: files[4],
				concurrence:  5,
				blockSize:    0,
			},
			sha256: sha256s[4],
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := make([]MfOption, 0)

			if tt.mf.Mf.url != "" {
				opts = append(opts, WithUrl(tt.mf.Mf.url))
			}
			if tt.mf.Mf.req != nil {
				u, err := url.Parse(tt.mf.Mf.url)
				if err != nil {
					t.Error(err)
					return
				}
				tt.mf.Mf.req.URL = u
				tt.mf.Mf.req.Method = tt.mf.Mf.method
				if tt.mf.Mf.req.Method == "" {
					tt.mf.Mf.req.Method = "GET"
				}
				opts = append(opts, WithReq(tt.mf.Mf.req))
			}
			if tt.mf.Mf.savePath != "" {
				opts = append(opts, WithSavePath(tt.mf.Mf.savePath))
			}
			if tt.mf.Mf.saveFileName != "" {
				opts = append(opts, WithSaveFileName(tt.mf.Mf.saveFileName))
			}
			if tt.mf.Mf.concurrence > 0 {
				opts = append(opts, WithConcurrence(tt.mf.Mf.concurrence))
			}
			if tt.mf.Mf.blockSize > 0 {
				opts = append(opts, WithBlockSize(tt.mf.Mf.blockSize))
			}
			mf, err := NewMf(opts...)
			if err != nil {
				t.Error(err)
				return
			}
			downloadFile := ""
			if downloadFile, err = mf.Resume(); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.mf.sha256 != "" {
				if err != nil {
					return
				}
				sha256Str, err := mtd.FileSha256(downloadFile)
				if err != nil {
					t.Error(err)
					return
				}
				if sha256Str != tt.mf.sha256 {
					t.Errorf("sha256Str = %v, want %v", sha256Str, tt.mf.sha256)
				}
			}
		})
	}
}
