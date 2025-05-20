package project

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/module"
	"github.com/agclqq/prow-framework/prowjob/command"
	strings2 "github.com/agclqq/prow-framework/strings"
)

var (
	frameworkModuleName = module.GetFrameworkName()
	moduleName          string
)

type Project struct {
}

func (a Project) GetCommand() string {
	return "init:project"
}
func (a Project) Usage() string {
	return `Usage of init:project:
  init:project [-t=base|full] []
`
}
func (a Project) Handle(ctx *prowjob.Context) {
	mn, err := module.GetName()
	if err != nil {
		fmt.Println(err)
		return
	}
	moduleName = mn

	//创建目录结构
	err = a.createDirs(ctx)
	if err != nil {
		fmt.Println("create dir error: ", err)
		return
	}

	//创建相关的文件
	err = a.createFiles(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//更新项目依赖
	err = a.tidy(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//格式化文件
	err = a.formatFiles(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 创建目录结构
func (a Project) createDirs(ctx *prowjob.Context) error {
	dirPaths := []string{
		"app/event",
		"app/service",
		"boot",
		"cmd/grpc",
		"cmd/httpd",
		"cmd/job",
		"config",
		"domain/demo/service",
		"infra/acl",
		"infra/db",
		"infra/queue",
		"res/views",
		"ui/console/command",
		"ui/console/register",
		"ui/grpc",
		"ui/grpc/controller",
		"ui/grpc/pb/demo",
		"ui/grpc/router",
		"ui/http",
		"ui/http/controller",
		"ui/http/middleware",
		"ui/http/response",
		"ui/http/router",
	}
	for _, dirPath := range dirPaths {
		dirPath = filepath.Clean(dirPath)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return err
		}
	}
	return nil
}

func (a Project) createFiles(ctx *prowjob.Context) error {
	err := a.createReadMeFiles(ctx)
	if err != nil {
		return err
	}
	err = a.createMakefileFiles(ctx)
	if err != nil {
		return err
	}
	err = a.createEnvFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createConfigFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createBootFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createHttpFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createCommandFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createEventFiles(ctx)
	if err != nil {
		return err
	}

	//创建demoagg相关的文件
	err = a.createDomainDemoFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createViewFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createGrpcFiles(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (a Project) createMakefileFiles(ctx *prowjob.Context) error {
	data := command.MakefileData{
		Vars: []string{
			"GO_BIN_DIR=bin",
		},
		MakefileRules: []command.MakefileRule{
			{
				Target:       ".PHONY",
				Dependencies: "http-local",
			},
			{
				Target:       "httpd-local",
				Dependencies: "",
				Commands: []string{
					"go build -o $(GO_BIN_DIR)/$@/$@ -v cmd/$@/*.go",
					"mkdir -p $(GO_BIN_DIR)/$@",
					"cp .env* $(GO_BIN_DIR)/$@",
				},
			},
			{
				Target:       "httpd-linux",
				Dependencies: "",
				Commands: []string{
					"CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GO_BIN_DIR)/httpd/httpd -v cmd/httpd/*.go",
					"mkdir -p $(GO_BIN_DIR)/$@",
					"cp .env* $(GO_BIN_DIR)/$@/",
				},
			},
			{
				Target:       "clean",
				Dependencies: "",
				Commands: []string{
					"rm -rf bin",
				},
			},
		},
	}
	return command.CreateTemplateFile("Makefile", command.MakefileTemplate, data)
}
func (a Project) createReadMeFiles(ctx *prowjob.Context) error {
	data := command.TextLineData{
		TextLines: []string{
			"# 简介",
			"# 项目介绍",
		},
	}
	return command.CreateTemplateFile("Readme.md", command.TextLines, data)
}

// 创建.env文件
func (a Project) createEnvFiles(ctx *prowjob.Context) error {
	data := command.TextLineData{
		TextLines: []string{
			"APP_ENV = dev",
			"HTTP_SERVER_PORT = 8080",
			"GRPC_SERVER_PORT = 8081",

			"#加密",
			"EASY_CRYPT_TYPE = aes/ECB/PKCS7/Hex",
			"EASY_CRYPT_KEY = " + a.createRand(16),
			"EASY_CRYPT_IV=" + a.createRand(16),

			"#log",
			"LOG_FILE	= log/app.log",
			"LOG_RETAIN	= 7",

			"#event",
			"EVENT_TEST = test",
			"EVENT_TEST_CAPACITY = 10",

			"#mysql",
			"DB_DEMO_HOST=127.0.0.1",
			"DB_DEMO_PORT=3306",
			"DB_DEMO_DB=demo",
			"DB_DEMO_DB_ALIAS=demo",
			"DB_DEMO_NAME=root",
			"DB_DEMO_PASS=",
			"DB_DEMO_DRIVER=mysql",
			"DB_DEMO_CHARSET=utf8mb4",
			"DB_DEMO_SQLLOG=logs/db/demo.log",
			"DB_DEMO_MAX_IDLE=10",
			"DB_DEMO_MAX_OPEN=10",
			"DB_DEMO_MAX_IDLE_TIME=25",
		},
	}

	return command.CreateTemplateFile(".env", command.TextLines, data)
}
func (a Project) createRand(length int) string {
	randStr := "abcdefghijklmnopqrstuvwxyz0123456789-"
	rs := strings2.GetRandomStr([]rune(randStr), length)
	return rs
}

// 创建config/app.go文件
func (a Project) createConfigFiles(ctx *prowjob.Context) error {
	err := a.createAppConfig(ctx)
	if err != nil {
		return err
	}
	err = a.createDbConfig(ctx)
	if err != nil {
		return err
	}
	err = a.createEventConfig(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (a Project) createAppConfig(ctx *prowjob.Context) error {
	data := command.TemplateData{}
	data.Imports = []command.ImportTemplate{{ImportName: frameworkModuleName + "/env/autoenv"}}
	confData := command.ConfTemplate{ConfName: "app"}
	confVars := make(command.Kv)
	confVars["appEnv"] = "APP_ENV"
	confVars["httpServerPort"] = "HTTP_SERVER_PORT"
	confVars["grpcServerPort"] = "GRPC_SERVER_PORT"
	confVars["easyCryptType"] = "EASY_CRYPT_TYPE"
	confVars["easyCryptKey"] = "EASY_CRYPT_KEY"
	confVars["easyCryptIv"] = "EASY_CRYPT_IV"
	confVars["logFile"] = "LOG_FILE"
	confVars["logRetain"] = "LOG_RETAIN"
	confData.Vars = confVars
	data.ConfData = confData
	return command.CreateTemplateFile("config/app.go", command.ConfigSingleMapTemplate, data)
}
func (a Project) createDbConfig(ctx *prowjob.Context) error {
	data := command.TemplateData{
		Imports: []command.ImportTemplate{
			{ImportName: frameworkModuleName + "/env/autoenv"},
		},
		ConfData: command.ConfTemplate{
			ConfName: "db",
			VarsM: command.KvM{
				"demo": {
					"host":        "DB_DEMO_HOST",
					"port":        "DB_DEMO_PORT",
					"db":          "DB_DEMO_DB",
					"alias":       "DB_DEMO_DB_ALIAS",
					"user":        "DB_DEMO_NAME",
					"password":    "DB_DEMO_PASS",
					"driver":      "DB_DEMO_DRIVER",
					"charset":     "DB_DEMO_CHARSET",
					"log":         "DB_DEMO_SQLLOG",
					"maxIdle":     "DB_DEMO_MAX_IDLE",
					"maxOpen":     "DB_DEMO_MAX_OPEN",
					"maxLife":     "DB_DEMO_MAX_LIFE",
					"maxIdleTime": "DB_DEMO_MAX_IDLE_TIME",
				},
			},
		},
	}
	return command.CreateTemplateFile("config/db.go", command.ConfigDoubleMapTemplate, data)
}
func (a Project) createEventConfig(ctx *prowjob.Context) error {
	data := command.TemplateData{}
	data.Imports = []command.ImportTemplate{{ImportName: frameworkModuleName + "/env/autoenv"}}
	confData := command.ConfTemplate{ConfName: "event"}
	confVars := make(command.KvM)
	confVars["test"] = map[string]string{
		"name":     "EVENT_TEST",
		"capacity": "EVENT_TEST_CAPACITY",
	}
	confData.VarsM = confVars
	data.ConfData = confData
	return command.CreateTemplateFile("config/event.go", command.ConfigDoubleMapTemplate, data)
}

func (a Project) createBootFiles(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "boot",
		Imports: []command.ImportTemplate{
			{ImportName: "fmt"},
			{ImportName: "strings"},
			{ImportName: "sync"},
			{ImportName: "gorm.io/gorm"},
			{ImportName: "github.com/spf13/cast"},
			{ImportName: moduleName + "/config"},
			{ImportName: moduleName + "/app/event/register"},
			{ImportName: frameworkModuleName + "/db"},
			{ImportName: frameworkModuleName + "/event"},
		},
		Vars: []string{
			"onceDbW = &sync.Once{}",
			"onceEvent = &sync.Once{}",
			"dbW    *gorm.DB",
		},
		Funcs: []command.FuncTemplate{
			{FuncName: "GetDbW", Params: "", ResultType: "*gorm.DB", FuncBody: "onceDbW.Do(func() {\n\t\t// 初始化mysql读写连接\n\t\tdbW = db.GetConn(\"demo\", config.GetDb(\"demo\"))\n\t})\n\treturn dbW"},
			{FuncName: "StartEvent", FuncBody: "onceEvent.Do(func() {\n\tvs:=[]string{}\n\tfor _, v := range config.GetAllEvent() {\n\t\terr := event.InitEvent(v[\"name\"], cast.ToInt(v[\"capacity\"]))\n\t\tif err != nil {\n\t\t\tfmt.Println(\"init event error \", err)\n\tcontinue\n\t}\n\tvs=append(vs, v[\"name\"])\n\t}\n\tregister.Register()\n\tevent.Run()\n\tfmt.Printf(\"%d events are running. \\n detail:%s\",len(vs),strings.Join(vs, \",\"))})"},
		},
	}
	return command.CreateTemplateFile("boot/boot.go", command.CommonTemplate, data)
}

func (a Project) createHttpFiles(ctx *prowjob.Context) error {
	err := a.createHttpd(ctx)
	if err != nil {
		return err
	}
	err = a.createHttpController(ctx)
	if err != nil {
		return err
	}
	err = a.createHttpControllerResp(ctx)
	if err != nil {
		return err
	}
	err = a.createHttpRouter(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 创建cmd/httpd/main.go文件
func (a Project) createHttpd(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "main",
		Imports: []command.ImportTemplate{
			{ImportName: "context"},
			{ImportName: "errors"},
			{ImportName: "fmt"},
			{ImportName: "net/http"},
			{Alias: "_", ImportName: "net/http/pprof"},
			{ImportName: "os"},
			{ImportName: "os/signal"},
			{ImportName: "sync"},
			{ImportName: "time"},
			{ImportName: moduleName + "/boot"},
			{ImportName: moduleName + "/ui/http/router"},
			{ImportName: moduleName + "/config"},
			{ImportName: frameworkModuleName + "/logger"},
			{Alias: "_", ImportName: frameworkModuleName + "/validator"},
			{ImportName: "github.com/spf13/cast"},
			{ImportName: "github.com/gin-gonic/gin"},
		},
	}
	return command.CreateTemplateFile("cmd/httpd/main.go", command.HttpdTemplate, data)
}

// 创建ui/http/controller/demo.go文件
func (a Project) createHttpController(ctx *prowjob.Context) error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Index", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller index demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Show", Params: "ctx *gin.Context", ResultType: "", FuncBody: "agg:=demo.NewAgg()\n\tresp,err:=agg.GetTest(ctx,1)\n\tif err!=nil{\n\t\tresponse.Failure(ctx,response.SERVER_ERROR,err.Error())\n\t\treturn \n\t}\n\tresponse.Success(ctx,resp)"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Update", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller update demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Store", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller store demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Destroy", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller destroy demo\"})"})

	data := command.TemplateData{
		PackageName: "controller",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/gin-gonic/gin"}, {ImportName: "net/http"}, {ImportName: moduleName + "/domain/demo"}, {ImportName: moduleName + "/ui/http/response"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{Name: typeName}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("ui/http/controller/demo.go", command.CommonTemplate, data)
}
func (a Project) createHttpControllerResp(ctx *prowjob.Context) error {
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{FuncName: "Index", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller index demo\"})"})
	funcs = append(funcs, command.FuncTemplate{FuncName: "Show", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller show demo\"})"})
	funcs = append(funcs, command.FuncTemplate{FuncName: "Update", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller update demo\"})"})
	funcs = append(funcs, command.FuncTemplate{FuncName: "Store", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller store demo\"})"})
	funcs = append(funcs, command.FuncTemplate{FuncName: "Destroy", Params: "ctx *gin.Context", ResultType: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller destroy demo\"})"})

	data := command.TemplateData{
		PackageName: "response",
		Imports:     []command.ImportTemplate{{ImportName: "context"}},
		Consts:      nil,
		Vars: []string{
			"SUCCESS           = resData{0, \"成功\"}",
			"INVALID_PARAMETER = resData{1001, \"无效的参数\"}",
			"INVALID_IDENTITY  = resData{1002, \"身份验证失败\"}",
			"PERMISSION_DENIED = resData{1003, \"没有权限\"}",
			"RESULT_EMPTY      = resData{1004, \"没有结果\"}",
			"SERVER_ERROR      = resData{3001, \"服务器错误\"}",
			"TIMEOUT           = resData{4001, \"访问超时\"}",
		},
		Types: []command.TypeTemplate{{Name: "resData", Fields: []string{"no  int", "msg string"}}, {Name: "Response", Fields: []string{"Code int         `json:\"code\"`", "Msg  string      `json:\"msg\"`", "Data interface{} `json:\"data,omitempty\"`"}}},
		Funcs: []command.FuncTemplate{
			{FuncName: "Success", Params: "ctx context.Context, data ...interface{}", ResultType: "Response", FuncBody: "return getResponse(ctx, SUCCESS, \"\", data...)"},
			{FuncName: "SuccessWithMsg", Params: "ctx context.Context, msg string, data ...interface{}", ResultType: "Response", FuncBody: "return getResponse(ctx, SUCCESS, msg, data...)"},
			{FuncName: "Failure", Params: "ctx context.Context, resErr resData, msg string, data ...interface{}", ResultType: "Response", FuncBody: "return getResponse(ctx, resErr, msg, data...)"},
			{FuncName: "getResponse", Params: "ctx context.Context, resErr resData, msg string, data ...interface{}", ResultType: "Response", FuncBody: "resMsg := resErr.msg\n\n\tif msg != \"\" {\n\t\tresMsg += \",\" + msg\n\t}\n\tif len(data) > 0 {\n\t\treturn Response{\n\t\t\tCode: resErr.no,\n\t\t\tMsg:  resMsg,\n\t\t\tData: data[0],\n\t\t}\n\t}\n\treturn Response{\n\t\tCode: resErr.no,\n\t\tMsg:  resMsg,\n\t}"},
		},
	}
	return command.CreateTemplateFile("ui/http/response/response.go", command.CommonTemplate, data)
}

// 创建ui/http/router/router.go文件
func (a Project) createHttpRouter(ctx *prowjob.Context) error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	var funcs []command.FuncTemplate

	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Register", Params: "eng *gin.Engine", ResultType: "", FuncBody: "apiGroup := eng.Group(\"/api\")\n\t{\n\trouter.ApiResource(apiGroup, \"/demo\", &controller.Demo{})\n\t}"})

	data := command.TemplateData{
		PackageName: "router",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/gin-gonic/gin"}, {ImportName: moduleName + "/ui/http/controller"}, {ImportName: frameworkModuleName + "/http/restful/router"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("ui/http/router/router.go", command.CommonTemplate, data)
}

func (a Project) createGrpcFiles(ctx *prowjob.Context) error {
	err := a.protocGenGoVersion()
	if err != nil {
		return err
	}
	err = a.protocGenGoGrpcVersion()
	if err != nil {
		return err
	}

	err = a.createGrpc(ctx)
	if err != nil {
		return err
	}
	err = a.createGrpcController(ctx)
	if err != nil {
		return err
	}
	err = a.createGrpcProto()
	if err != nil {
		return err
	}
	err = a.buildProto()
	if err != nil {
		return err
	}
	err = a.createGrpcRouter()
	if err != nil {
		return err
	}
	return nil
}

func (a Project) protocGenGoVersion() error {
	//protocGenGoVersion := "v1.28.1"
	_, err := execCommand("protoc-gen-go", "--version")
	//if err != nil {
	//	return err
	//}
	//version := strings.Split(rs, " ")
	//if len(version) < 2 {
	//	return errors.New("want protoc-gen-go version,but got " + rs)
	//}
	//
	//if compareVersions(version[1], protocGenGoVersion) < 0 {
	//	return errors.New("it is recommended that protoc-gen-go be at least " + protocGenGoVersion)
	//}
	return err
}

func (a Project) protocGenGoGrpcVersion() error {
	//protocGenGoGrpcVersion := "1.3.0"
	_, err := execCommand("protoc-gen-go-grpc", "--version")
	return err
	//if err != nil {
	//	return err
	//}
	//version := strings.Split(rs, " ")
	//if len(version) < 2 {
	//	return errors.New("want protoc-gen-go-grpc version,but got " + rs)
	//}
	//
	//if compareVersions(version[1], protocGenGoGrpcVersion) < 0 {
	//	return errors.New("it is recommended that protoc-gen-go-grpc be at least " + protocGenGoGrpcVersion)
	//}
	//return nil
}

func (a Project) buildProto() error {
	rs, err := execCommand("protoc", "-I=ui/grpc/pb", "--go_out=ui/grpc/pb", "--go-grpc_out=ui/grpc/pb", "ui/grpc/pb/demo.proto")
	if err != nil {
		fmt.Println(rs)
		return err
	}
	return nil
}

func execCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = cmd.Stdout
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	rs, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return string(rs), err
	}
	return string(rs), nil
}

func compareVersions(v1, v2 string) int {
	// 定义版本号的正则表达式模式
	versionPattern := regexp.MustCompile(`(\d+\.\d+\.\d+)`)

	// 从字符串中提取版本号
	extractVersion := func(version string) []int {
		matches := versionPattern.FindStringSubmatch(version)
		if len(matches) != 2 {
			return nil
		}

		// 将版本号分割为整数切片
		nums := make([]int, 3)
		for i, v := range regexp.MustCompile(`\d+`).FindAllString(matches[1], -1) {
			num, _ := strconv.Atoi(v)
			nums[i] = num
		}
		return nums
	}

	// 提取版本号
	v1Nums := extractVersion(v1)
	v2Nums := extractVersion(v2)

	// 比较版本号大小
	for i := 0; i < 3; i++ {
		if v1Nums[i] < v2Nums[i] {
			return -1 // v1 < v2
		} else if v1Nums[i] > v2Nums[i] {
			return 1 // v1 > v2
		}
	}
	return 0 // 版本号相等
}

// 创建cmd/grpc/main.go文件
func (a Project) createGrpc(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "",
		Imports: []command.ImportTemplate{
			{ImportName: "errors"},
			{ImportName: "fmt"},
			{ImportName: "net"},
			{ImportName: "net/http"},
			{ImportName: "os"},
			{ImportName: "os/signal"},
			{ImportName: "sync"},
			{ImportName: "time"},
			{ImportName: "google.golang.org/grpc"},
			{ImportName: "google.golang.org/grpc/keepalive"},
			{Alias: "grpcRouter", ImportName: moduleName + "/ui/grpc/router"},
			{ImportName: moduleName + "/config"},
		},
	}

	return command.CreateTemplateFile("cmd/grpc/main.go", command.GrpcTemplate, data)
}

// 创建ui/grpc/controller/demo.go文件
func (a Project) createGrpcController(ctx *prowjob.Context) error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Foo", Params: "context.Context, *demo.DemoReq", ResultType: "(*demo.DemoRes, error)", FuncBody: "return &demo.DemoRes{Name: \"Hello World\"}, nil"})

	data := command.TemplateData{
		PackageName: "controller",
		Imports:     []command.ImportTemplate{{ImportName: "context"}, {ImportName: "google.golang.org/grpc"}, {ImportName: moduleName + "/ui/grpc/pb/demo"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{Name: typeName, Fields: []string{"Server *grpc.Server", "demo.UnimplementedDemoServer"}}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("ui/grpc/controller/demo.go", command.CommonTemplate, data)
}

// 创建ui/grpc/pb/demo.proto文件
func (a Project) createGrpcProto() error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Foo", Params: "DemoReq", ResultType: "DemoRes", FuncBody: ""})

	var msg []command.MessageTemplate
	msg = append(msg, command.MessageTemplate{MsgName: "DemoReq", Fields: []string{"string name = 1;"}})
	msg = append(msg, command.MessageTemplate{MsgName: "DemoRes", Fields: []string{"string name = 1;"}})

	data := command.TemplateData{
		PackageName: "/demo",
		Imports:     nil,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{Name: typeName}},
		Funcs:       funcs,
		ProtoPkg:    typeName,
		Services:    []command.ServiceTemplate{{Name: typeName, Funcs: funcs}},
		Messages:    msg,
	}
	return command.CreateTemplateFile("ui/grpc/pb/demo.proto", command.ProtoTemplate, data)
}

// 创建ui/grpc/router/router.go文件
func (a Project) createGrpcRouter() error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Register", Params: "s *grpc.Server", ResultType: "", FuncBody: "demo.RegisterDemoServer(s, &controller.Demo{Server: s})"})

	data := command.TemplateData{
		PackageName: "router",
		Imports:     []command.ImportTemplate{{ImportName: moduleName + "/ui/grpc/controller"}, {ImportName: moduleName + "/ui/grpc/pb/demo"}, {ImportName: "google.golang.org/grpc"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("ui/grpc/router/router.go", command.CommonTemplate, data)
}

func (a Project) createCommandFiles(ctx *prowjob.Context) error {
	err := a.createCommandCmd(ctx)
	if err != nil {
		return err
	}
	err = a.createCommand(ctx)
	if err != nil {
		return err
	}
	err = a.createCommandRegister(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (a Project) createCommandCmd(ctx *prowjob.Context) error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "main", Params: "", ResultType: "", FuncBody: "jobEng := prowjob.New()\nprowjobreg.Register(jobEng)\nregister.Register(jobEng)\njobEng.Run()"})

	data := command.TemplateData{
		PackageName: "main",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}, {Alias: "prowjobreg", ImportName: frameworkModuleName + "/prowjob/register"}, {ImportName: moduleName + "/ui/console/register"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("cmd/job/main.go", command.CommonTemplate, data)
}

// 创建ui/console/command/demo.go文件
func (a Project) createCommand(ctx *prowjob.Context) error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "GetCommand", Params: "", ResultType: "string", FuncBody: "return \"command:demo\""})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Usage", Params: "", ResultType: "string", FuncBody: "return `Usage of command:demo:\n\t  command:demo\n\t`"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Handle", Params: "ctx *prowjob.Context", ResultType: "", FuncBody: "fmt.Println(\"this is command demo\")"})

	data := command.TemplateData{
		PackageName: "command",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}, {ImportName: "fmt"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{Name: typeName}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("ui/console/command/demo.go", command.CommonTemplate, data)
}

// 创建ui/console/command/register.go文件
func (a Project) createCommandRegister(ctx *prowjob.Context) error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	mn, err := module.GetName()
	if err != nil {
		return err
	}
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Register", Params: "eng *prowjob.CommandEngine", ResultType: "", FuncBody: "eng.Add(&command.Demo{})"})
	data := command.TemplateData{
		PackageName: "register",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}, {ImportName: mn + "/ui/console/command"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("ui/console/register/register.go", command.CommonTemplate, data)
}

func (a Project) createEventFiles(ctx *prowjob.Context) error {
	err := a.createEventRegister(ctx)
	if err != nil {
		return err
	}
	err = a.createEventDemoFiles(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (a Project) createEventRegister(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "register",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prow-framework/event"}, {ImportName: moduleName + "/app/event"}},
		Types:       []command.TypeTemplate{{Name: "Demo"}},
		Funcs: []command.FuncTemplate{
			{
				FuncName: "Register",
				FuncBody: "event.Register(&event.Demo{})",
			},
		},
	}
	return command.CreateTemplateFile("app/event/register/register.go", command.CommonTemplate, data)
}
func (a Project) createEventDemoFiles(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "event",
		Imports:     []command.ImportTemplate{{ImportName: "context"}, {ImportName: "fmt"}},
		Types:       []command.TypeTemplate{{Name: "Demo"}},
		Funcs: []command.FuncTemplate{
			{
				Receiver:     "d",
				ReceiverType: "*" + "Demo",
				FuncName:     "ListenName",
				Params:       "",
				ResultType:   "string",
				FuncBody:     "return \"test\"",
			},
			{
				Receiver:     "d",
				ReceiverType: "*" + "Demo",
				FuncName:     "Concurrence",
				Params:       "",
				ResultType:   "int64",
				FuncBody:     "return 1",
			},
			{
				Receiver:     "d",
				ReceiverType: "*" + "Demo",
				FuncName:     "Handle",
				Params:       "ctx context.Context, data []byte",
				ResultType:   "",
				FuncBody:     "fmt.Println(\"demo event\", string(data))",
			},
		},
	}
	return command.CreateTemplateFile("app/event/demo.go", command.CommonTemplate, data)
}

// 创建domain示例相关的文件
func (a Project) createDomainDemoFiles(ctx *prowjob.Context) error {
	//err := a.createDemoAggRoot(ctx)
	//if err != nil {
	//	return err
	//}
	err := a.createDemoAggService(ctx)
	if err != nil {
		return err
	}
	err = a.createDemoAgg(ctx)
	if err != nil {
		return err
	}
	err = a.createDemoAggRepo(ctx)
	if err != nil {
		return err
	}
	err = a.createDemoAggEntity(ctx)
	if err != nil {
		return err
	}
	err = a.createDemoAggVo(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 创建domain/aggdemo/root.go文件
func (a Project) createDemoAggRoot(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "demo",
		Imports:     []command.ImportTemplate{{ImportName: "context"}},
		Consts:      nil,
		Vars:        nil,
		Interfaces:  []command.InterTemplate{{Name: "Demo", Methods: []string{"GetTest(ctx context.Context,id int) (*EntityA, error)"}}},
	}
	return command.CreateTemplateFile("domain/demo/root.go", command.CommonTemplate, data)
}

func (a Project) createDemoAgg(ctx *prowjob.Context) error {
	typeName := "DemoAgg"
	receiver := strings.ToLower(typeName[0:1])
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{FuncName: "NewAgg", Params: "", ResultType: "*DemoAgg", FuncBody: "return &DemoAgg{}"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "GetTest", Params: "ctx context.Context,id int", ResultType: "(*EntityA,error)", FuncBody: "repo:=&DemoRepo{}\n\treturn repo.Select(ctx,&EntityA{Id: id}),nil"})

	data := command.TemplateData{
		PackageName: "demo",
		Imports:     []command.ImportTemplate{{ImportName: "context"}},
		Consts:      nil,
		Vars:        []string{"_ Demo=(*DemoAgg)(nil)"},
		Types:       []command.TypeTemplate{{Name: typeName}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("domain/demo/agg.go", command.CommonTemplate, data)
}

// 创建domain/aggdemo/entity.go文件
func (a Project) createDemoAggEntity(ctx *prowjob.Context) error {
	typeName := "EntityA"

	data := command.TemplateData{
		PackageName: "demo",
		Imports:     nil,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{Name: typeName, Fields: []string{"Id int    `json:\"id\" gorm:\"primary_key;auto_increment\"`", "Name string `json:\"name omitempty\"`", "Status int `json:\"status\"`"}}},
	}

	return command.CreateTemplateFile("domain/demo/entity.go", command.CommonTemplate, data)
}
func (a Project) createDemoAggRepo(ctx *prowjob.Context) error {
	typeName := "Repo"
	typeNameImpl := typeName + "Impl"
	entityName := "EntityA"
	tableName := "demo_table"
	receiver := strings.ToLower(typeName[0:1])
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{FuncName: "New" + typeName, Params: "db *gorm.DB", ResultType: typeName, FuncBody: "return &" + typeNameImpl + "{table: \"" + tableName + "\", db: db, R: repo.NewRepo[" + entityName + "](db, \"" + tableName + "\")}"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "TableName", Params: "", ResultType: "string", FuncBody: "return \"" + entityName + "\""})
	data := command.TemplateData{
		PackageName: "demo",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prow-framework/db/repo"}, {ImportName: "gorm.io/gorm"}},
		Consts:      nil,
		Vars:        []string{"_ Repo = (*" + typeNameImpl + ")(nil)"},
		Interfaces:  []command.InterTemplate{{Name: typeName, Methods: []string{"repo.R[" + entityName + "]", "TableName() string"}}},
		Types: []command.TypeTemplate{{Name: typeNameImpl, Fields: []string{"table string",
			"db    *gorm.DB",
			"repo.R[" + entityName + "]"}}},
		Funcs: funcs,
	}

	return command.CreateTemplateFile("domain/demo/repo.go", command.CommonTemplate, data)
}
func (a Project) createDemoAggService(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "demo",
	}
	return command.CreateTemplateFile("domain/demo/service.go", command.CommonTemplate, data)
}

// 创建domain/demoagg/vo.go文件
func (a Project) createDemoAggVo(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "demo",
		Imports:     nil,
		Consts: []string{
			"Var_a =iota",
			"Var_b",
			"Var_c",
		},
	}
	return command.CreateTemplateFile("domain/demo/vo.go", command.CommonTemplate, data)
}

func (a Project) createViewFiles(ctx *prowjob.Context) error {
	data := command.TemplateData{
		ViewTemplateDefine: "{{define \"index/index.tmpl\"}}",
		ViewTemplateEnd:    "{{end}}",
	}
	return command.CreateTemplateFile("res/views/index/index.tmpl", command.TmplTemplate, data)
}

func (a Project) tidy(ctx *prowjob.Context) error {
	s, err := execCommand("go", "mod", "tidy")
	if err != nil {
		fmt.Println(s)
		return err
	}
	return nil
}

func (a Project) formatFiles(ctx *prowjob.Context) error {
	s, err := execCommand("gofmt", "-w", "-s", ".")
	if err != nil {
		fmt.Println(s)
		return err
	}
	s, err = execCommand("goimports", "-local", moduleName, "-w", ".")
	if err != nil {
		fmt.Println(s)
		return err
	}
	return nil
}
