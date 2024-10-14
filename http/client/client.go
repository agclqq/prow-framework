package client

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"strings"
	"time"

	"github.com/agclqq/prow-framework/skywalking"

	ginPluginsHttp "github.com/SkyAPM/go2sky/plugins/http"
)

type ByteBody []byte
type DiyResponse struct {
	*http.Response
	ByteBody
}

type HttpClient struct {
	Req *http.Request
	cli *http.Client
}

var (
	defaultClient = NewClient()
	globalClient  *HttpClient
	transport     = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,  // 播号超时时间，默认为操作系统的超时时间，linux默认为3分钟
			KeepAlive: 65 * time.Second, // 保持连接的时长，默认为15秒，如果设置为负数表示不保持连接
		}).DialContext,
		ForceAttemptHTTP2: true,

		IdleConnTimeout:       90 * time.Second,                      // 空闲连接失效时间，DefaultTransport set 90
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, //#nosec G402 -- 不校验服务端证书
		TLSHandshakeTimeout:   10 * time.Second,                      // DefaultTransport set 10
		MaxIdleConns:          100,                                   // 最大空闲连接个数，DefaultTransport set 100
		MaxIdleConnsPerHost:   10,                                    // 每个客户端最大空闲连接个数，DefaultMaxIdleConnsPerHost = 2
		MaxConnsPerHost:       100,                                   // 每个客户端可以持有的最大连接个数,默认为0
		ExpectContinueTimeout: 1 * time.Second,                       //  DefaultTransport set 1
	}
)

type Option func(hc *HttpClient)

func GetClient() *HttpClient {
	globalClient = NewClient()
	return globalClient
}

func NewClient(opts ...Option) *HttpClient { //
	cli := &HttpClient{
		cli: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second, // DefaultTransport set 0
		},
	}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}
func WithTimeout(dur time.Duration) Option {
	return func(hc *HttpClient) {
		hc.cli.Timeout = dur
	}
}
func WithCookieJar(jar http.CookieJar) Option {
	return func(hc *HttpClient) {
		hc.cli.Jar = jar
	}
}
func WithTrace() Option {
	return func(hc *HttpClient) {
		hc.WithTrace()
	}
}

func (h *HttpClient) WithTrace() *HttpClient {
	tracer := skywalking.GetTracer()
	if tracer == nil {
		return h
	}
	transValue := reflect.ValueOf(h.cli.Transport)
	if _, ok := h.cli.Transport.(*http.Transport); !ok {
		return h
	}
	if transValue.Type() != reflect.TypeOf(&http.Transport{}) {
		//没有作更详细的判断，因为ginPluginsHttp.NewClient会改写client的transport为ginPluginsHttp.transport
		return h
	}
	client, err := ginPluginsHttp.NewClient(tracer, ginPluginsHttp.WithClient(h.cli))
	if err != nil {
		return h
	}
	h.cli = client
	return h
}

func (h *HttpClient) Original() *http.Client {
	return h.cli
}
func (h *HttpClient) SetTimeOut(t time.Duration) {
	h.cli.Timeout = t
}
func (h *HttpClient) SetHeader(key, val string) {
	h.Req.Header.Set(key, val)
}
func (h *HttpClient) SetCookie(c *http.Cookie) {
	h.Req.AddCookie(c)
}
func (h *HttpClient) SetCookieJar(jar *cookiejar.Jar) {
	h.cli.Jar = jar
}
func (h *HttpClient) GetCookieJar() http.CookieJar {
	return h.cli.Jar
}

func (h *HttpClient) SetBearerAuth(token string) error {
	if h.Req == nil {
		return errors.New("you should build a request at first")
	}
	h.Req.Header.Set("Authorization", "Bearer "+token)
	return nil
}
func (h *HttpClient) SetBasicAuth(user, password string) error {
	if h.Req == nil {
		return errors.New("you should build a request at first")
	}
	h.Req.SetBasicAuth(user, password)
	return nil
}
func Request(method string, url string, body io.Reader) (*HttpClient, error) {
	return defaultClient.Request(method, url, body)
}
func RequestWithCtx(ctx context.Context, method string, url string, body io.Reader) (*HttpClient, error) {
	return defaultClient.RequestWithCtx(ctx, method, url, body)
}
func (h *HttpClient) Request(method string, url string, body io.Reader) (*HttpClient, error) {
	ctx := context.Background()
	return h.RequestWithCtx(ctx, method, url, body)
}
func (h *HttpClient) RequestWithCtx(ctx context.Context, method string, url string, body io.Reader) (*HttpClient, error) {
	method = strings.ToUpper(method)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	h.Req = req
	return h, err
}
func (h *HttpClient) Send() (*http.Response, error) {
	return h.cli.Do(h.Req)
}

// SendAutoClose 自动读取body的内容，并关闭，返回的结构体为DiyResponse
func (h *HttpClient) SendAutoClose() (*DiyResponse, error) {
	res, err := h.Send()
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close() // #nosec G104
	response := &DiyResponse{
		Response: res,
		ByteBody: body,
	}
	return response, nil
}

//func SendWithRetry(ctx context.Context, url string, method string,
//	headers http.Header, body io.Reader, rty *retry.Retry) (*DiyResponse, error) {
//	var res *DiyResponse
//	var err error
//	var cli = GetClient()
//	defer func() {
//		if err != nil {
//			logger.Error(ctx, "http request failed, err: %v. url: %v, method: %v, header: %v, body: %s",
//				err, url, method, headers, body)
//		}
//	}()
//
//	rty.Run(func(step uint) {
//		cli, err = cli.RequestWithCtx(ctx, method, url, body)
//		if err != nil {
//			rty.Cancel()
//			logger.ErrorWithTrace(ctx, err)
//			return
//		}
//		if headers != nil {
//			cli.Req.Header = headers
//		}
//		res, err = cli.SendAutoClose()
//		if err != nil {
//			logger.ErrorfWithTrace(ctx, "request failed, url: %v, method: %v, err: %v, retry now, step is %d",
//				url, method, err, step)
//			return
//		}
//		rty.Cancel()
//	})
//	return res, err
//}
//
//func SendOnlyWithRetry(ctx context.Context, cli *HttpClient, rty *retry.Retry) (*DiyResponse, error) {
//	var res *DiyResponse
//	var err error
//	rty.Run(func(step uint) {
//		res, err = cli.SendAutoClose()
//		if err != nil {
//			logger.ErrorfWithTrace(ctx, "request failed, err: %v, retry now, step is %d", err, step)
//			return
//		}
//		rty.Cancel()
//	})
//	return res, err
//}
