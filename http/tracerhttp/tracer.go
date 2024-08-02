package tracerhttp

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/SkyAPM/go2sky"
	skyhttp "github.com/SkyAPM/go2sky/plugins/http"
)

var poolClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 播号超时时间，默认为操作系统的超时时间，linux默认为3分钟
			KeepAlive: 30 * time.Second, // 保持连接的时长，默认为15秒，如果设置为负数表示不保持连接
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,              // DefaultTransport set 100
		IdleConnTimeout:       90 * time.Second, // DefaultTransport set 90
		TLSHandshakeTimeout:   10 * time.Second, // DefaultTransport set 10
		MaxIdleConnsPerHost:   2,                // DefaultTransport set 2
		ExpectContinueTimeout: 1 * time.Second,  //  DefaultTransport set 1
	},
	Timeout: 30, // DefaultTransport set 0
}
var DefaultClient *http.Client

func PoolOpt() skyhttp.ClientOption {
	return skyhttp.WithClient(poolClient)
}
func GetClient(tracer *go2sky.Tracer, opt ...skyhttp.ClientOption) (*http.Client, error) {
	if len(opt) == 0 {
		opt = append(opt, PoolOpt())
	}
	//http.GetInfo()
	DefaultClient, err := skyhttp.NewClient(tracer, opt...)
	return DefaultClient, err
}
func Get(url string) {

}
func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {

	return
}
