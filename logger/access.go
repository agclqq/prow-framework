package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	lt "github.com/agclqq/prow-framework/times"

	"github.com/gin-gonic/gin"
)

type Config struct {
	RequestBody bool
	gin.LoggerConfig
}

func WithConfig(conf Config) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		//
	}

	out := conf.Output
	if out == nil {
		//out = DefaultWriter
	}

	notlogged := conf.SkipPaths

	//isTerm := true

	//if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
	//	(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
	//	isTerm = false
	//}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var body []byte
		if conf.RequestBody {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			if conf.RequestBody {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}

func AccessStdConfig(e *gin.Engine) Config {
	return Config{
		RequestBody: true,
		LoggerConfig: gin.LoggerConfig{
			Formatter: AccessLogJsonFormatter(e),
			Output:    os.Stdout,
			SkipPaths: nil,
		},
	}
}

func AccessLogConfig(e *gin.Engine, file string, retain int) Config {
	if retain <= 0 {
		retain = DefaultRetain
	}
	writer, err := RotateDailyLog(file, uint(retain))
	if err != nil {
		log.Fatal(err)
	}
	return Config{
		RequestBody: true,
		LoggerConfig: gin.LoggerConfig{
			Formatter: AccessLogJsonFormatter(e),
			Output:    writer,
			SkipPaths: nil,
		},
	}
}

func AccessLogFormatter() func(gin.LogFormatterParams) string {
	return func(p gin.LogFormatterParams) string {
		//nginx $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"
		//this  $remote_addr [$time_local] $http_method  $status "$request"  "$http_referer" "$http_user_agent"
		return fmt.Sprintf("%s [%s] %s \"%d\" %s \"%s\" \"%s\"\n",
			strings.Split(p.Request.RemoteAddr, ":")[0],
			p.TimeStamp.Format(lt.FormatDatetimeMicro),
			p.Method,
			p.StatusCode,
			p.Request.RequestURI,
			p.Request.Header.Get("Referer"),
			p.Request.Header.Get("User-Agent"),
		)
	}
}

type accessLogJson struct {
	RemoteAddr     string `json:"remote_addr"`
	TimeLocal      string `json:"time_local"`
	HttpMethod     string `json:"http_method"`
	Status         int    `json:"status"`
	Request        string `json:"request"`
	HttpReferer    string `json:"http_referer"`
	HttpUserAgent  string `json:"http_user_agent"`
	HttpXForwarded string `json:"http_x_forwarded"`
	BodyBytesSent  int    `json:"body_bytes_sent"`
	RequestBody    string `json:"request_body"`
	RequestTime    string `json:"request_time"`
}

func AccessLogJsonFormatter(e *gin.Engine) func(params gin.LogFormatterParams) string {
	return func(p gin.LogFormatterParams) string {
		//nginx $remote_addr - $remote_user [$time_local]                "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"
		//this  $remote_addr                 $time_local   $http_method  $request   $status                  "$http_referer" "$http_user_agent"
		var body []byte
		if strings.HasPrefix(p.Request.Header.Get("Content-Type"), "multipart/form-data") {
			if len(p.Request.PostForm) == 0 {
				p.Request.ParseMultipartForm(e.MaxMultipartMemory) // #nosec G104
			}
			body = []byte(p.Request.PostForm.Encode())
		} else {
			bodyReader := p.Request.Body
			body, _ = io.ReadAll(bodyReader)
			p.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		oneMB := 1 << 10
		if len(body) > oneMB {
			body = body[:oneMB]
		}

		alj := accessLogJson{
			RemoteAddr:     strings.Split(p.Request.RemoteAddr, ":")[0],
			TimeLocal:      p.TimeStamp.Format(lt.FormatDatetimeMicro),
			HttpMethod:     p.Method,
			Status:         p.StatusCode,
			Request:        p.Request.RequestURI,
			HttpReferer:    p.Request.Header.Get("Referer"),
			HttpUserAgent:  p.Request.Header.Get("User-Agent"),
			HttpXForwarded: p.Request.Header.Get("X-Forwarded-For"),
			RequestBody:    string(body),
			BodyBytesSent:  p.BodySize,
			RequestTime:    fmt.Sprintf("%v", p.Latency),
		}
		jsonLog, err := json.Marshal(alj)
		if err != nil {
			log.Print(err)
		}
		str := string(jsonLog) + "\n"
		return str
	}
}
