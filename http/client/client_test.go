package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"

	"github.com/agclqq/prow-framework/skywalking"
)

func TestHttpClient_Request(t *testing.T) {
	//t.Skip("此测试函数不会被执行")
	cli := NewClient()
	_, err := cli.Request("get", "https://httpbin.org/", nil)
	if err != nil {
		return
	}
	res, err := cli.SendAutoClose()
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != 200 {
		t.Error("请求失败")
		return
	}
}

func setTracer() {
	re, err := reporter.NewLogReporter()
	if err != nil {
		fmt.Printf("new reporter error %v \n", err)
		return
	}
	defer re.Close()
	tracer, err := go2sky.NewTracer("WithTrace_test", go2sky.WithReporter(re))
	if err != nil {
		return
	}
	skywalking.SetTracer(tracer)
}
func TestHttpClient_GetCookieJar(t *testing.T) {
	//t.Skip("此测试函数不会被执行")
	ctx := context.Background()
	var err error
	var data = make(url.Values)
	data.Set("principal", "testName")
	data.Set("password", "testPassword")

	cli := GetClient()

	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	cli.SetCookieJar(jar)

	cli, err = cli.RequestWithCtx(ctx, "GET", "https://httpbin.org/c/log_out", nil)
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	if err != nil {
		return
	}
	rs, err := cli.SendAutoClose()
	if err != nil {
		t.Error(err)
		return
	}
	if "" == rs.Header.Get("X-Harbor-Csrf-Token") {
		t.Error("want header key :X-Harbor-Csrf-Token")
		return
	}

	cli, err = cli.RequestWithCtx(ctx, "POST", "https://httpbin.org/c/login", strings.NewReader(data.Encode()))
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	cli.SetHeader("X-Harbor-Csrf-Token", rs.Header.Get("X-Harbor-Csrf-Token"))
	cli.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return
	}
	rs, err = cli.SendAutoClose()
	if err != nil {
		return
	}
	if rs.StatusCode != http.StatusOK {
		t.Error(string(rs.ByteBody))
	}
	got := cli.GetCookieJar()
	cookies := got.Cookies(cli.Req.URL)
	fmt.Println(cookies)
	if len(cookies) == 0 {
		t.Errorf("GetCookieJar() = %v, len=%d ,want len 1", got, len(cookies))
	}
	cli, err = cli.RequestWithCtx(ctx, "GET", "https://httpbin.org/api/v2.0/users/current", nil)
	if err != nil {
		return
	}
	cli.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	rs, err = cli.SendAutoClose()
	if err != nil {
		t.Error(err)
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(rs.ByteBody, &m)
	if err != nil {
		t.Error(err)
		return
	}
	if v, ok := m["username"]; !ok {
		t.Error("get null")
	} else {
		t.Log(v)
	}
}

func TestGetClient(t *testing.T) {
	lc := sync.Mutex{}
	wg := sync.WaitGroup{}
	m := make(map[string]int)
	re, err := regexp.Compile(`=(\d)`)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		wg.Add(2)
		ii := i
		go func() {
			defer wg.Done()
			cli1 := GetClient()
			cli1, err := cli1.Request("get", "https://httpbin.org/?s=1"+strconv.Itoa(ii), nil)
			if err != nil {
				t.Error(err)
				return
			}
			res1, err := cli1.SendAutoClose()
			if err != nil {
				t.Error(err)
				return
			}

			parts := re.FindStringSubmatch(res1.Request.URL.RawQuery)
			lc.Lock()
			m[parts[1]]++
			lc.Unlock()
			t.Logf("res1 %d,%s", ii, res1.Request.URL.RequestURI())
		}()
		go func() {
			defer wg.Done()
			cli2 := GetClient()
			cli2, err := cli2.Request("get", "https://httpbin.org/?s=2"+strconv.Itoa(ii), nil)
			if err != nil {
				t.Error(err)
				return
			}
			res2, err := cli2.SendAutoClose()
			if err != nil {
				t.Error(err)
				return
			}
			parts := re.FindStringSubmatch(res2.Request.URL.RawQuery)
			lc.Lock()
			m[parts[1]]++
			lc.Unlock()
			t.Logf("res2 %d, %s", ii, res2.Request.URL.RequestURI())
		}()
	}

	wg.Wait()
	t.Logf("%+v", m)
}

func TestGetClient1(t *testing.T) {
	cli1 := GetClient()
	var err error
	go func() {
		defer func() {
			if errAny := recover(); errAny != nil {
				t.Error(errAny)
			}
		}()
		cli1, err = cli1.Request("get", "https://httpbin.org/?s=1", nil)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(3 * time.Second)
		res1, err := cli1.SendAutoClose()
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("res1 %s", res1.Request.URL.RequestURI())
	}()
	go func() {

		time.Sleep(1 * time.Second)
		cli1, err := cli1.Request("get", "https://httpbin.org/?s=2", nil)
		time.Sleep(3 * time.Second)
		if err != nil {
			t.Error(err)
			return
		}
		res2, err := cli1.SendAutoClose()
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("res2 %s", res2.Request.URL.RequestURI())
	}()
	time.Sleep(10 * time.Second)
}

func TestHttpClient_WithTrace(t *testing.T) {
	setTracer()

	for i := 0; i < 1; i++ {
		span, _, err := skywalking.AddTrace(context.Background(), skywalking.GetTracer(), "test-span")
		if err != nil {
			fmt.Println(err)
		}
		cli1 := GetClient().WithTrace()
		cli1, err = cli1.Request("get", "https://httpbin.org/?s=1", nil)
		if err != nil {
			return
		}
		cli1.Send()
		cli2 := GetClient().WithTrace()
		cli2, err = cli2.Request("get", "https://httpbin.org/?s=2", nil)
		if err != nil {
			return
		}
		cli2.Send()
		span.End()
	}
}

func TestHttpClient_ShouldBodyClose(t *testing.T) {
	const numRequests = 100

	// 创建自定义的 Transport，允许长连接保持时间为1秒
	transport := &http.Transport{
		IdleConnTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: transport}
	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.Get("https://httpbin.org")
			if err != nil {
				t.Error(err)
				return
			}
			//defer resp.Body.Close()

			// 在这里处理响应
		}()
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
}
func TestHttpClient_Original_retry(t *testing.T) {
	for i := 0; i < 10; i++ {
		TestHttpClient_ShouldBodyClose(t)
	}
}

func TestHttpClient_AlternateBody(t *testing.T) {
	client := http.Client{}
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		res, err := client.Get("https://httpbin.org")
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(2 * time.Second)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf(string(body))
	}()
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		res, err := client.Get("https://httpbin.org")
		if err != nil {
			t.Error(err)
			return
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf(string(body))
	}()
	wg.Wait()
}

func TestNewClient(t *testing.T) {
	setTracer()
	cli1 := NewClient(WithTimeout(10))
	cli2 := NewClient(WithTimeout(2 * time.Second))
	cli3 := NewClient(WithTimeout(3*time.Second), WithTrace())
	cli4 := NewClient(WithTimeout(4*time.Second), WithTrace())
	fmt.Printf("1:%p\n2:%p\n3:%p\n4:%p\n", cli1.cli.Transport, cli2.cli.Transport, cli3.cli.Transport, cli4.cli.Transport)
	if cli1 == cli2 {
		t.Error("Expectations are not equal, but reality is equal")
		return
	}
	if cli1.cli.Transport != cli2.cli.Transport {
		t.Error("Expectations is equal, but reality are not equal")
		return
	}
	if cli1.cli.Transport == cli3.cli.Transport {
		t.Error("Expectations are not equal, but reality is equal")
		return
	}
	if cli4.cli.Transport == cli3.cli.Transport {
		t.Error("Expectations are not equal, but reality is equal")
		return
	}
	_, err := cli1.cli.Get("https://httpbin.org")
	if err == nil {
		t.Error("Expect a timeout error, but it does not") //期望是出错的，因为过期时间太短了，完不成一次响应
		return
	}

	_, err = cli2.cli.Get("https://httpbin.org")
	if err != nil {
		t.Error(err)
		return
	}
}
