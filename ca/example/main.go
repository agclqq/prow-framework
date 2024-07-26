package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/agclqq/prow-framework/ca/example/ca"
	"github.com/agclqq/prow-framework/ca/example/client"
	"github.com/agclqq/prow-framework/ca/example/server"
)

func main() {
	os.Remove("ca.key")
	os.Remove("ca.crt")
	os.Remove("caSvr.key")
	os.Remove("caSvr.crt")
	os.Remove("svr.key")
	os.Remove("svr.crt")
	os.Remove("cli.key")
	os.Remove("cli.crt")
	defer func() {
		os.Remove("ca.key")
		os.Remove("ca.crt")
		os.Remove("caSvr.key")
		os.Remove("caSvr.crt")
		os.Remove("svr.key")
		os.Remove("svr.crt")
		os.Remove("cli.key")
		os.Remove("cli.crt")
	}()
	//启动CA服务
	fmt.Println("init ca cert and caSvr")
	go func() {
		err := ca.Svr("caSvr.crt", "caSvr.key")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	time.Sleep(2 * time.Second)

	//获取CA证书
	fmt.Println("get ca cert")

	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := cli.Get("https://127.0.0.1:8080/cert")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return
	}
	caCert, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	//启动服务端
	go func() {
		fmt.Println("init ca cert and caSvr")
		err = server.Svr(caCert, "svr.crt", "svr.key")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	time.Sleep(2 * time.Second)
	//启动客户端
	fmt.Println("request svr")
	err = client.Cli(caCert, "cli.key", "cli.crt")
	if err != nil {
		fmt.Println(err)
		return
	}
}
