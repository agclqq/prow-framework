package cookiejar

import (
	"net/http"
	"net/url"
)

type Supper struct {
	cookies []*http.Cookie
}

func (jar *Supper) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *Supper) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

func (jar *Supper) SetCookie(cookie *http.Cookie) {
	jar.cookies = append(jar.cookies, cookie)
}
