package make

import (
	"errors"
	"fmt"
	"strings"

	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/args"
	"github.com/agclqq/prow-framework/codefmt"
	"github.com/agclqq/prow-framework/db"
	"github.com/agclqq/prow-framework/prowjob/command"
	strings2 "github.com/agclqq/prow-framework/strings"

	"gorm.io/gorm"
)

const DEFAULT_MODEL_DIR = "infra/repo/"

type Model struct {
}
type mysqlTableStruct struct {
	Field   string
	Type    string
	Key     string
	Comment string
}

var mysqlTableStructs []mysqlTableStruct

func (a Model) GetCommand() string {
	return "make:model"
}
func (a Model) Usage() string {
	return `Usage of make:model:
 make:model dbType dns tableName [--alias] [--path]
   alias affects the package name and defaults is db name
   if the 'path' is not given,the default is '` + DEFAULT_MODEL_DIR + `${driver}/${db}'
`
}
func (a Model) Handle(ctx *prowjob.Context) {
	if len(ctx.Param) < 3 {
		fmt.Printf("%s \n %s", command.ERROR_PARAM_NUM, a.Usage())
		return
	}
	pparam := args.TidyParmaWithPrefix(ctx.Param)
	if _, ok := pparam["h"]; ok {
		fmt.Println(a.Usage())
		return
	}
	dbType := ctx.Param[0]
	dsn := ctx.Param[1]
	tableName := ctx.Param[2]

	dbConfDsn, err := db.DsnDecode(dbType, dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbConf := make(map[string]string)
	dbConf["driver"] = dbConfDsn.Type
	dbConf["host"] = dbConfDsn.Host
	dbConf["port"] = dbConfDsn.Port
	dbConf["user"] = dbConfDsn.User
	dbConf["password"] = dbConfDsn.Password
	dbConf["db"] = dbConfDsn.Dbname
	dbConf["charset"] = dbConfDsn.Charset
	dbConf["alias"] = pparam["alias"]
	dbConf["timezone"] = dbConfDsn.TimeZone
	dbConf["alias"] = pparam["alias"]
	if dbConf["alias"] == "" {
		dbConf["alias"] = dbConf["db"]
	}
	dbConf["table"] = tableName
	path := strings.Trim(pparam["path"], "/") + "/"
	fullPath, err := getModelFullPath(dbConf, path)
	if err != nil {
		fmt.Printf(command.NO_DB_CONFIG+" \n", dbType)
		return
	}
	if !CheckOverwrite(fullPath) { //Abort if the file exists and the user does not allow it to be overwritten
		return
	}

	descs := getTableStruct(dbType, dbConf, tableName)

	err = createModelFile(dbConf, fullPath, descs)
	if err != nil {
		fmt.Println(err)
		return
	}
	ff, err := formatFile(fullPath)
	fmt.Println(ff)
	if err != nil {
		fmt.Println(err)
	}
}

func getModelFullPath(dbConf map[string]string, path string) (string, error) {
	if dbConf == nil {
		return "", errors.New(fmt.Sprintf(command.NO_DB_CONFIG, dbConf))
	}
	if path == "" {
		path = DEFAULT_MODEL_DIR + dbConf["driver"] + "/" + dbConf["alias"] + "/"
	}
	f := path + dbConf["table"] + ".go"
	return f, nil
}
func getTableStruct(configName string, config map[string]string, tableName string) [][]string {
	conn := db.GetConn(configName, config)
	if conn == nil {
		conn.Name()
		fmt.Printf(command.NO_DB_CONFIG, configName)
		return nil
	}
	var desc [][]string
	switch conn.Name() {
	case "sqlite":
		//fallthrough
	case "mysql":
		desc = getDescFromMysql(conn, tableName)
	}
	return desc
}

func getDescFromMysql(conn *gorm.DB, tableName string) [][]string {
	conn.Table(tableName).Raw("show full columns from " + tableName).Scan(&mysqlTableStructs)
	var rs [][]string
	for _, v := range mysqlTableStructs {
		field := v.Field
		if v.Field == "id" && v.Key == "PRI" {
			field = "ID"
		}
		r := append([]string{}, strings2.SnakeToBigCamel(field), resolveTypes(v.Type), createTags(v), createComment(v.Comment))
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
	if strings.Contains(fieldType, "json") {
		return "datatypes.JSON"
	}
	return "" //未支持的类型
}
func resolveKeys(fieldKey string) string {
	if fieldKey == "PRI" {
		return "`gorm:\"primaryKey\"`"
	}
	return ""
}
func createTags(v mysqlTableStruct) string {
	keyTag := ""
	jsonTag := ""
	if resolveKeys(v.Key) != "" {
		keyTag = "gorm:\"primaryKey\""
	}

	if v.Field == "create_time" && v.Type == "datetime" {
		keyTag = "gorm:\"autoCreateTime\""
	}
	if v.Field == "update_time" && (v.Type == "datetime" || v.Type == "timestamp") {
		keyTag = "gorm:\"autoUpdateTime\""
	}

	jsonTag = "json:\"" + v.Field + "\""
	s := append([]string{}, jsonTag, keyTag)
	return "`" + strings.TrimSpace(strings.Join(s, " ")) + "`"
}
func createComment(comment string) string {
	if comment == "" {
		return ""
	}
	comment = strings.ReplaceAll(comment, "\n", "\t//")
	return " // " + comment
}

func createModelFile(conf map[string]string, modelPath string, desc [][]string) error {
	if modelPath == "" {
		return errors.New("model path cannot be empty")
	}
	packageName := conf["alias"]
	typeName := strings2.ToUpFirst(strings2.SnakeToBigCamel(conf["table"]))
	receiver := strings2.ToLowFirst(typeName[0:1])
	receiverType := "*" + typeName
	importStat := make(map[string]bool)
	var importList []command.ImportTemplate
	var content []string
	for _, v := range desc {
		if v[1] == "time.Time" && !importStat["time"] {
			importList = append(importList, command.ImportTemplate{ImportName: "time"})
			importStat["time"] = true
		}
		if v[1] == "datatypes.JSON" && !importStat["datatypes.JSON"] {
			importList = append(importList, command.ImportTemplate{ImportName: "gorm.io/datatypes"})
			importStat["datatypes.JSON"] = true
		}
		row := strings.Join(v, "\t")
		content = append(content, "\t"+row)
	}
	return creatFile(packageName, typeName, receiver, receiverType, conf["table"], modelPath, content, importList)
}

func creatFile(packageName, typeName, receiver, receiverType, tableName, modelPath string, content []string, importList []command.ImportTemplate) error {
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{FuncName: "New" + typeName, Params: "", ResultType: receiverType, FuncBody: "return &" + typeName + "{}"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "TableName", Params: "", ResultType: "string", FuncBody: "return \"" + tableName + "\""})

	data := command.TemplateData{
		PackageName: packageName,
		Imports:     importList,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{Name: typeName, Fields: content}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile(modelPath, command.CommonTemplate, data)
}

func formatFile(path string) (string, error) {
	rs := ""
	goFmt, err := codefmt.GoFmt()
	if err != nil {
		return rs, err
	}
	rs += "\n" + goFmt
	imports, err := codefmt.GoImports(path)
	if err != nil {
		return rs, err
	}
	rs += "\n" + imports
	return rs, nil
}
