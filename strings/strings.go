package strings

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"unsafe"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!@#$%^&*()")

func ToUpFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	f := []rune(str)
	if f[0] >= 97 && f[0] <= 122 {
		f[0] = f[0] - 32
		return string(f)
	}
	return str
}
func ToLowFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	f := []rune(str)
	if f[0] >= 65 && f[0] <= 90 {
		f[0] = f[0] + 32
		return string(f)
	}
	return str
}

func SnakeToBigCamel(str string) string {
	if str == "" {
		return str
	}
	strs := strings.Split(str, "_")
	var strBui []byte
	for _, v := range strs {
		strBui = append(strBui, []byte(ToUpFirst(v))...)
	}
	return string(strBui)
}

func TrimAllTrimSuffix(oldS, suffix string) string {
	newS := strings.TrimSuffix(oldS, suffix)
	if newS != oldS {
		return TrimAllTrimSuffix(newS, suffix)
	}
	return newS
}

// 以下两方法从gin框中摘取

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	// #nosec G103
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b)) // #nosec G103
}

func TokenGenerator(n int) string {
	token := make([]rune, n)
	for i := range token {
		b, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		token[i] = letters[b.Int64()]
	}
	return string(token)
}

// GetRandomString 生成随机字符串
func GetRandomString(n int) string {
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return GetRandomStr(letters, n)
}
func GetRandomStr(letters []rune, n int) string {
	token := make([]rune, n)
	for i := range token {
		b, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		token[i] = letters[b.Int64()]
	}
	return string(token)
}

func GetCopyName(n int) string {
	randomString := GetRandomString(n)
	randomString = strings.ToLower(randomString)
	return "_c_" + randomString
}

// ConvToStr 转string
func ConvToStr(num interface{}) string {
	return fmt.Sprintf("%v", num)
}

// GoFormat 类似python传入字典格式化方法
func GoFormat(format string, p map[string]interface{}) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

func GetStringLen(s string) int {
	count := 0
	for k := range s {
		_ = k
		count++
	}
	return count
}

func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

// ValidateString 字母,数字,下划线,中划线,点，且不能以点和中划线开头
func ValidateString(str string) bool {
	pattern := "^[^-.][a-zA-Z0-9_\\-.]+$"
	regx := regexp.MustCompile(pattern)
	return regx.MatchString(str)
}

// ValidateString2 字母,数字,中划线, 且不能以中划线开头
func ValidateString2(str string) bool {
	pattern := "^[^-_][a-zA-Z0-9\\-.]+$"
	regx := regexp.MustCompile(pattern)
	return regx.MatchString(str)
}
