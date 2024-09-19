package encoding

import (
	"encoding/json"
	"os"

	"github.com/tidwall/gjson"
)

// AnyToJsonStrList 结构体转 json string
func AnyToJsonStrList(obj interface{}) string {
	var jsonBytes, err = json.Marshal(obj)
	if err != nil {
		empty := []struct{}{}
		jsonBytes, _ := json.Marshal(empty)
		return string(jsonBytes)
	}
	jsonStr := string(jsonBytes)
	if jsonStr == "null" {
		empty := []struct{}{}
		jsonBytes, _ := json.Marshal(empty)
		return string(jsonBytes)
	}
	return jsonStr
}

// ToJsonStr 结构体转 json string
func ToJsonStr(obj interface{}) string {
	var jsonStr, err = json.Marshal(obj)
	//var jsonStr, err = json.MarshalIndent(obj, "", "    ")
	if err != nil {
		os.Exit(-1)
	}
	return string(jsonStr)
}

// StrToJsonStr golang记录的string字符串日志里无法格式化，这里帮助格式化输出
func StrToJsonStr(s string) string {
	defer func() {
		if err := recover(); err != nil {
			os.Exit(-1)
		}
	}()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}
	return ToJsonStr(m)
}

// GetJsonField 不反序列化，只取某个字段
func GetJsonField(jsonStr, key string) string {
	return gjson.Get(jsonStr, key).String()
}
