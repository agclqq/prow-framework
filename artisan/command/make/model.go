package make

//
//import (
//	"errors"
//	"fmt"
//	"os"
//	"strings"
//
//	"github.com/agclqq/prow-framework/args"
//	"github.com/agclqq/prow-framework/artisan"
//	file2 "github.com/agclqq/prow-framework/file"
//	"github.com/agclqq/prow-framework/repository/manager"
//	strings2 "github.com/agclqq/prow-framework/strings"
//
//	"gorm.io/gorm"
//)
//
//const DEFAULT_MODEL_DIR = "infrastructure/repository/"
//
//type Model struct {
//}
//type mysqlTableStruct struct {
//	Field   string
//	Type    string
//	Key     string
//	Comment string
//}
//
//var mysqlTableStructs []mysqlTableStruct
//
//func (a Model) GetCommand() string {
//	return "make:model"
//}
//func (a Model) Usage() string {
//	return `Usage of make:model:
//  make:model configName tableName [path]
//    if the 'path' is not given,the default is 'infrastructure/repository/${driver}/${db}'
//`
//}
//func (a Model) Handle(ctx *artisan.Context) {
//	if len(ctx.Param) < 3 {
//		fmt.Printf("%s \n %s", artisan.ERROR_PARAM_NUM, a.Usage())
//	}
//	pparam := args.TidyParmaWithPrefix(ctx.Param)
//	if _, ok := pparam["h"]; ok {
//		fmt.Println(a.Usage())
//		return
//	}
//	connConfig := ctx.Param[1]
//	tableName := ctx.Param[2]
//	path := ""
//	if len(ctx.Param) == 4 {
//		path = ctx.Param[3]
//	}
//	f, err := getModelFullPath(connConfig, tableName, path)
//	if err != nil {
//		fmt.Printf(artisan.NO_DB_CONFIG+" \n", connConfig)
//	}
//	if !checkOverwrite(f) { //Abort if the file exists and the user does not allow it to be overwritten
//		return
//	}
//
//	descs := getTableStruct(connConfig, tableName)
//
//	err = createModelFile(connConfig, tableName, f, descs)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//func getModelFullPath(dbConf map[string]string, tableName, path string) (string, error) {
//	if dbConf == nil {
//		return "", errors.New(fmt.Sprintf(artisan.NO_DB_CONFIG, dbConf))
//	}
//	if path == "" {
//		path = DEFAULT_MODEL_DIR + dbConf["driver"] + "/" + dbConf["alias"] + "/"
//	}
//	f := path + tableName + ".go"
//	return f, nil
//}
//func getTableStruct(configName, tableName string) [][]string {
//	conn := manager.GetConn(configName)
//	if conn == nil {
//		conn.Name()
//		fmt.Printf(artisan.NO_DB_CONFIG, configName)
//		return nil
//	}
//	var desc [][]string
//	switch conn.Name() {
//	case "sqlite":
//		//fallthrough
//	case "mysql":
//		desc = getDescFromMysql(conn, tableName)
//	}
//	return desc
//}
//
//func getDescFromMysql(conn *gorm.DB, tableName string) [][]string {
//	conn.Table(tableName).Raw("show full columns from " + tableName).Scan(&mysqlTableStructs)
//	var rs [][]string
//	for _, v := range mysqlTableStructs {
//		field := v.Field
//		if v.Field == "id" && v.Key == "PRI" {
//			field = "ID"
//		}
//		r := append([]string{}, strings2.SnakeToBigCamel(field), resolveTypes(v.Type), createTags(v), createComment(v.Comment))
//		rs = append(rs, r)
//	}
//	return rs
//}
//func resolveTypes(fieldType string) string {
//	if strings.HasPrefix(fieldType, "bit") {
//		return "bool"
//	}
//	if strings.HasPrefix(fieldType, "bigint") {
//		if strings.Contains(fieldType, "unsigned") {
//			return "uint64"
//		}
//		return "int64"
//	}
//	if strings.HasPrefix(fieldType, "int") || strings.HasPrefix(fieldType, "mediumint") {
//		if strings.Contains(fieldType, "unsigned") {
//			return "uint"
//		}
//		return "int"
//	}
//	if strings.HasPrefix(fieldType, "smallint") {
//		if strings.Contains(fieldType, "unsigned") {
//			return "uint16"
//		}
//		return "int16"
//	}
//	if strings.HasPrefix(fieldType, "tinyint") {
//		if strings.Contains(fieldType, "unsigned") {
//			return "uint8"
//		}
//		return "int8"
//	}
//	if strings.HasPrefix(fieldType, "float") {
//		return "float32"
//	}
//	if strings.HasPrefix(fieldType, "double") || strings.HasPrefix(fieldType, "decimal") {
//		return "float64"
//	}
//	if strings.Contains(fieldType, "char") ||
//		strings.Contains(fieldType, "text") ||
//		strings.Contains(fieldType, "blob") {
//		return "string"
//	}
//	if strings.Contains(fieldType, "date") || strings.Contains(fieldType, "time") { // 包含了fieldType == "datetime" || fieldType == "timestamp"
//		return "time.Time"
//	}
//	if strings.Contains(fieldType, "json") {
//		return "datatypes.JSON"
//	}
//	return "" //未支持的类型
//}
//func resolveKeys(fieldKey string) string {
//	if fieldKey == "PRI" {
//		return "`gorm:\"primaryKey\"`"
//	}
//	return ""
//}
//func createTags(v mysqlTableStruct) string {
//	keyTag := ""
//	jsonTag := ""
//	if resolveKeys(v.Key) != "" {
//		keyTag = "gorm:\"primaryKey\""
//	}
//
//	if v.Field == "create_time" && v.Type == "datetime" {
//		keyTag = "gorm:\"autoCreateTime\""
//	}
//	if v.Field == "update_time" && (v.Type == "datetime" || v.Type == "timestamp") {
//		keyTag = "gorm:\"autoUpdateTime\""
//	}
//
//	jsonTag = "json:\"" + v.Field + "\""
//	s := append([]string{}, jsonTag, keyTag)
//	return "`" + strings.TrimSpace(strings.Join(s, " ")) + "`"
//}
//func createComment(comment string) string {
//	if comment == "" {
//		return ""
//	}
//	return " // " + comment
//}
//
//func createModelFile(connConfig, tableName, modelPath string, desc [][]string) error {
//	conf := config.GetDb(connConfig)
//	if conf == nil {
//		return errors.New(fmt.Sprintf(artisan.NO_DB_CONFIG, connConfig))
//	}
//	if modelPath == "" {
//		return errors.New("model path cannot be empty")
//	}
//	pgName := conf["alias"]
//	tyName := strings2.ToUpFirst(strings2.SnakeToBigCamel(tableName))
//	importStat := make(map[string]bool)
//	var importList []string
//	var content []string
//	for _, v := range desc {
//		if v[1] == "time.Time" && !importStat["time"] {
//			importList = append(importList, "\"time\"")
//			importStat["time"] = true
//		}
//		if v[1] == "datatypes.JSON" && !importStat["datatypes.JSON"] {
//			importList = append(importList, "\"gorm.io/datatypes\"")
//			importStat["datatypes.JSON"] = true
//		}
//		row := strings.Join(v, "\t")
//		content = append(content, "\t"+row)
//	}
//	importStr := ""
//	if len(importList) > 0 {
//		importStr = `
//import (
//` + strings.Join(importList, "\n") + `
//)
//`
//	}
//	contentStr := strings.Join(content, "\n")
//	src := `package ` + pgName + `
//` + importStr + `
//type ` + tyName + ` struct {
//` + contentStr + `
//}
//func (t *` + tyName + `) TableName() string {
//	return "` + tableName + `"
//}
//`
//	if err := file2.MakeDirByFile(modelPath); err != nil {
//		return err
//	}
//	if err := os.WriteFile(modelPath, []byte(src), 0666); err != nil {
//		return err
//	}
//	return nil
//}
