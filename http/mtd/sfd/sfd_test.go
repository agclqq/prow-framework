package sfd

import (
	"net/http"
	url "net/url"
	"testing"

	"github.com/agclqq/prow-framework/http/mtd"
)

type testSf struct {
	Sf     //inherit
	sha256 string
}

func TestSf_Download(t *testing.T) {
	var tests = []struct {
		name    string
		mf      testSf
		wantErr bool
	}{
		{name: "t1_dynamicText_url",
			mf: testSf{
				Sf: Sf{
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
		{name: "t1_dynamicText_req", mf: testSf{
			Sf: Sf{
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
		{name: "t2_bin_url", mf: testSf{
			Sf: Sf{
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
		{name: "t2_bin_req", mf: testSf{
			Sf: Sf{
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
		{name: "t3_img_url", mf: testSf{
			Sf: Sf{
				url:          "https://httpbin.org/cdn-static/luban/nextjs/images/default-large-files/1713205013445-1E217CCAC7AFA986.png",
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
		{name: "t3_img_req", mf: testSf{
			Sf: Sf{
				url:          "https://httpbin.org/cdn-static/luban/nextjs/images/default-large-files/1713205013445-1E217CCAC7AFA986.png",
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
			opts := make([]SfOption, 0)

			if tt.mf.Sf.url != "" {
				opts = append(opts, WithUrl(tt.mf.Sf.url))
			}
			if tt.mf.Sf.req != nil {
				u, err := url.Parse(tt.mf.Sf.url)
				if err != nil {
					t.Error(err)
					return
				}
				tt.mf.Sf.req.URL = u
				tt.mf.Sf.req.Method = tt.mf.Sf.method
				if tt.mf.Sf.req.Method == "" {
					tt.mf.Sf.req.Method = "GET"
				}
				opts = append(opts, WithReq(tt.mf.Sf.req))
			}
			if tt.mf.Sf.savePath != "" {
				opts = append(opts, WithSavePath(tt.mf.Sf.savePath))
			}
			if tt.mf.Sf.saveFileName != "" {
				opts = append(opts, WithSaveFileName(tt.mf.Sf.saveFileName))
			}
			if tt.mf.Sf.concurrence > 0 {
				opts = append(opts, WithConcurrence(tt.mf.Sf.concurrence))
			}
			if tt.mf.Sf.blockSize > 0 {
				opts = append(opts, WithBlockSize(tt.mf.Sf.blockSize))
			}
			sf, err := NewSfd(opts...)
			if err != nil {
				t.Error(err)
				return
			}
			downloadFile := ""
			if downloadFile, err = sf.Download(); (err != nil) != tt.wantErr {
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
