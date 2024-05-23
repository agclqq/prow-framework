package mfd

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/agclqq/prow-framework/http/mtd"
)

type testMf struct {
	Mf
	sha256 string
}

func TestMf_Download(t *testing.T) {
	var tests = []struct {
		name    string
		mf      testMf
		wantErr bool
	}{
		{name: "t1_dynamicText_url",
			mf: testMf{
				Mf: Mf{
					url:          "http://127.0.0.1:8080/apiex/log",
					method:       "",
					header:       nil,
					req:          nil,
					savePath:     "/tmp",
					saveFileName: "url.log.log",
					concurrence:  5,
					blockSize:    0,
				},
				sha256: "6ce00649806afd2ec9f97017e5512b81d2ab107ae166695ee33e115f6a6f8b9a",
			},
			wantErr: false},
		{name: "t1_dynamicText_req", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/apiex/log",
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     "/tmp",
				saveFileName: "req.log.log",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "",
		}, wantErr: false},
		{name: "t2_bin_url", mf: testMf{
			Mf: Mf{
				url:          "https://cdn.stubdownloader.services.mozilla.com/builds/firefox-latest-ssl/zh-CN/osx/85d91034c0e7a65c7eab91e65088198c8822df1c98026bb12397c18ff5d808ca/Firefox%20125.0.3.dmg",
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     "/tmp",
				saveFileName: "url.firefox.dmg",
				concurrence:  -1,
				blockSize:    0,
			},
			sha256: "5c8b535e4d87baeb27463643fcf2b5d9f5d816299925b428413108dc0483e812",
		}, wantErr: false},
		{name: "t2_bin_req", mf: testMf{
			Mf: Mf{
				url:          "https://cdn.stubdownloader.services.mozilla.com/builds/firefox-latest-ssl/zh-CN/osx/85d91034c0e7a65c7eab91e65088198c8822df1c98026bb12397c18ff5d808ca/Firefox%20125.0.3.dmg",
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     "/tmp",
				saveFileName: "req.firefox.dmg",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "5c8b535e4d87baeb27463643fcf2b5d9f5d816299925b428413108dc0483e812",
		}, wantErr: false},
		{name: "t3_img_url", mf: testMf{
			Mf: Mf{
				url:          "https://whttpbin.org/cdn-static/luban/nextjs/images/default-large-files/1713205013445-1E217CCAC7AFA986.png",
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     "/tmp",
				saveFileName: "",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "c0d5098271682ea48ee9be6776857e6c2599afd58e7acad481431f2686af674c",
		}, wantErr: false},
		{name: "t3_img_req", mf: testMf{
			Mf: Mf{
				url:          "https://whttpbin.org/cdn-static/luban/nextjs/images/default-large-files/1713205013445-1E217CCAC7AFA986.png",
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     "/tmp",
				saveFileName: "req.img.png",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "c0d5098271682ea48ee9be6776857e6c2599afd58e7acad481431f2686af674c",
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
	}
}

func TestMf_Resume(t *testing.T) {
	var tests = []struct {
		name    string
		mf      testMf
		wantErr bool
	}{
		{name: "t1_resume_dynamicText_url",
			mf: testMf{
				Mf: Mf{
					url:          "http://127.0.0.1:8080/apiex/log",
					method:       "",
					header:       nil,
					req:          nil,
					savePath:     "/tmp",
					saveFileName: "resume.url.log.log",
					concurrence:  5,
					blockSize:    0,
				},
				sha256: "6ce00649806afd2ec9f97017e5512b81d2ab107ae166695ee33e115f6a6f8b9a",
			},
			wantErr: false},
		{name: "t1_resume_dynamicText_req", mf: testMf{
			Mf: Mf{
				url:          "http://127.0.0.1:8080/apiex/log",
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     "/tmp",
				saveFileName: "resume.req.log.log",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "",
		}, wantErr: false},
		{name: "t2_resume_bin_url", mf: testMf{
			Mf: Mf{
				url:          "https://cdn.stubdownloader.services.mozilla.com/builds/firefox-latest-ssl/zh-CN/osx/85d91034c0e7a65c7eab91e65088198c8822df1c98026bb12397c18ff5d808ca/Firefox%20125.0.3.dmg",
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     "/tmp",
				saveFileName: "resume.url.firefox.dmg",
				concurrence:  -1,
				blockSize:    0,
			},
			sha256: "5c8b535e4d87baeb27463643fcf2b5d9f5d816299925b428413108dc0483e812",
		}, wantErr: false},
		{name: "t2_resume_bin_req", mf: testMf{
			Mf: Mf{
				url:          "https://cdn.stubdownloader.services.mozilla.com/builds/firefox-latest-ssl/zh-CN/osx/85d91034c0e7a65c7eab91e65088198c8822df1c98026bb12397c18ff5d808ca/Firefox%20125.0.3.dmg",
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     "/tmp",
				saveFileName: "resume.re.firefox.dmg",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "5c8b535e4d87baeb27463643fcf2b5d9f5d816299925b428413108dc0483e812",
		}, wantErr: false},
		{name: "t3_resume_img_url", mf: testMf{
			Mf: Mf{
				url:          "https://whttpbin.org/cdn-static/luban/nextjs/images/default-large-files/1713205013445-1E217CCAC7AFA986.png",
				method:       "",
				header:       nil,
				req:          nil,
				savePath:     "/tmp",
				saveFileName: "resume.url.img.png",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "c0d5098271682ea48ee9be6776857e6c2599afd58e7acad481431f2686af674c",
		}, wantErr: false},
		{name: "t3_resume_img_req", mf: testMf{
			Mf: Mf{
				url:          "https://whttpbin.org/cdn-static/luban/nextjs/images/default-large-files/1713205013445-1E217CCAC7AFA986.png",
				method:       "",
				header:       nil,
				req:          &http.Request{},
				savePath:     "/tmp",
				saveFileName: "resume.req.img.png",
				concurrence:  5,
				blockSize:    0,
			},
			sha256: "c0d5098271682ea48ee9be6776857e6c2599afd58e7acad481431f2686af674c",
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
