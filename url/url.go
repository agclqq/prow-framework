package url

import (
	url2 "net/url"
	"strings"
)

func GetBasePath(url string) (string, string, error) {
	parse, err := url2.Parse(url)
	if err != nil {
		return "", "", err
	}
	//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	base := strings.Builder{}
	if parse.Scheme != "" {
		base.Write([]byte(parse.Scheme + ":"))
	}
	if parse.Host != "" {
		base.Write([]byte("//"))
		if parse.User != nil {
			base.Write([]byte(parse.User.String() + "@"))
		}
		base.Write([]byte(parse.Host))
	}
	return base.String(), strings.Trim(parse.Path, "/"), nil
}

func GetGitRepoBathPath(url string) (string, error) {
	parse, err := url2.Parse(url)
	if err != nil {
		return "", err
	}
	base := strings.Builder{}
	if parse.Host != "" {
		base.Write([]byte(parse.Host))
	}

	if parse.Path != "" {
		base.Write([]byte(parse.Path))
	}

	return base.String(), nil
}
