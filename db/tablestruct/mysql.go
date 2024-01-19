package tablestruct

import (
	"strings"

	strings2 "github.com/agclqq/prow-framework/strings"

	"gorm.io/gorm"
)

type Mysql struct {
}
type MysqlTableStruct struct {
	Field string
	Type  string
	Key   string
}

var MysqlTableStructs []MysqlTableStruct

func (m Mysql) GetDesc(conn *gorm.DB, tableName string) [][]string {
	conn.Table(tableName).Raw("desc " + tableName).Scan(&MysqlTableStructs)
	var rs [][]string
	for _, v := range MysqlTableStructs {
		field := v.Field
		if v.Field == "id" && v.Key == "PRI" {
			field = "ID"
		}
		r := append([]string{}, strings2.SnakeToBigCamel(field), resolveTypes(v.Type), createTags(v.Field, v.Key))
		rs = append(rs, r)
	}
	return rs
}
func resolveTypes(fieldType string) string {
	if strings.HasPrefix(fieldType, "bit") {
		return "bool"
	}
	if strings.HasPrefix(fieldType, "bigint") {
		if strings.Contains(fieldType, "unsigned") {
			return "uint64"
		}
		return "int64"
	}
	if strings.HasPrefix(fieldType, "int") || strings.HasPrefix(fieldType, "mediumint") {
		if strings.Contains(fieldType, "unsigned") {
			return "uint"
		}
		return "int"
	}
	if strings.HasPrefix(fieldType, "smallint") {
		if strings.Contains(fieldType, "unsigned") {
			return "uint16"
		}
		return "int16"
	}
	if strings.HasPrefix(fieldType, "tinyint") {
		if strings.Contains(fieldType, "unsigned") {
			return "uint8"
		}
		return "int8"
	}
	if strings.HasPrefix(fieldType, "float") {
		return "float32"
	}
	if strings.HasPrefix(fieldType, "double") || strings.HasPrefix(fieldType, "decimal") {
		return "float64"
	}
	if strings.Contains(fieldType, "char") ||
		strings.Contains(fieldType, "text") ||
		strings.Contains(fieldType, "blob") {
		return "string"
	}
	if strings.Contains(fieldType, "date") || strings.Contains(fieldType, "time") { // 包含了fieldType == "datetime" || fieldType == "timestamp"
		return "time.Time"
	}
	return "" //未支持的类型
}
func createTags(field, key string) string {
	keyTag := ""
	jsonTag := ""
	if resolveKeys(key) != "" {
		keyTag = "gorm:\"primaryKey\""
	}
	jsonTag = "json:\"" + field + "\""
	s := append([]string{}, keyTag, jsonTag)
	return "`" + strings.Join(s, " ") + "`"
}
func resolveKeys(fieldKey string) string {
	if fieldKey == "PRI" {
		return "`gorm:\"primaryKey\"`"
	}
	return ""
}
