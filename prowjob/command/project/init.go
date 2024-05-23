package project

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/module"
	"github.com/agclqq/prow-framework/prowjob/command"
	strings2 "github.com/agclqq/prow-framework/strings"
)

var (
	frameworkModuleName = module.GetNameWithoutErr()
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

	//格式化文件
	err = a.formatFiles(ctx)
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
}

// 创建目录结构
func (a Project) createDirs(ctx *prowjob.Context) error {
	dirPaths := []string{
		"app/console/command",
		"app/events",
		"app/grpc/controller",
		"app/grpc/pb",
		"app/grpc/router",
		"app/http/controller",
		"app/http/router",
		"app/middleware",
		"bootstrap",
		"cmd/httpd",
		"cmd/grpc",
		"config",
		"domain",
		"infra",
		"resource/views",
	}
	for _, dirPath := range dirPaths {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (a Project) createFiles(ctx *prowjob.Context) error {
	err := a.createEnv(ctx)
	if err != nil {
		return err
	}

	err = a.createConfig(ctx)
	if err != nil {
		return err
	}

	err = a.createHttpFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createGrpcFiles(ctx)
	if err != nil {
		return err
	}

	err = a.createCommandFiles(ctx)
	if err != nil {
		return err
	}

	//创建demoagg相关的文件
	err = a.createDomainDemoFiles(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 创建.env文件
func (a Project) createEnv(ctx *prowjob.Context) error {
	data := command.TemplateData{}
	envData := make([]command.EventTemplate, 0)
	envData = append(envData, command.EventTemplate{Key: "APP_ENV", Val: "dev"})
	envData = append(envData, command.EventTemplate{Key: "HTTP_SERVER_PORT", Val: "8080"})
	envData = append(envData, command.EventTemplate{Key: "GRPC_SERVER_PORT", Val: "8081"})

	envData = append(envData, command.EventTemplate{Val: "#加密", Type: "comment"})
	envData = append(envData, command.EventTemplate{Key: "EASY_CRYPT_TYPE", Val: "aes/ECB/PKCS7/Hex"})
	envData = append(envData, command.EventTemplate{Key: "EASY_CRYPT_KEY", Val: a.createRand(16)})
	envData = append(envData, command.EventTemplate{Key: "EASY_CRYPT_IV", Val: a.createRand(16)})

	envData = append(envData, command.EventTemplate{Val: "#log", Type: "comment"})
	envData = append(envData, command.EventTemplate{Key: "LOG_FILE", Val: "log/app.log"})
	envData = append(envData, command.EventTemplate{Key: "LOG_RETAIN", Val: "7"})

	data.Envs = envData
	return command.CreateTemplateFile(".env", command.EnvTemplate, data)
}
func (a Project) createRand(length int) string {
	randStr := "abcdefghijklmnopqrstuvwxyz0123456789-"
	rs := strings2.GetRandomStr([]rune(randStr), length)
	return rs
}

// 创建config/app.go文件
func (a Project) createConfig(ctx *prowjob.Context) error {
	data := command.TemplateData{}
	data.Imports = []command.ImportTemplate{{ImportName: frameworkModuleName + "/env"}}
	confData := make(map[string]string)
	confData["appEnv"] = "APP_ENV"
	confData["httpServerPort"] = "HTTP_SERVER_PORT"
	confData["grpcServerPort"] = "GRPC_SERVER_PORT"
	confData["easyCryptType"] = "EASY_CRYPT_TYPE"
	confData["easyCryptKey"] = "EASY_CRYPT_KEY"
	confData["easyCryptIv"] = "EASY_CRYPT_IV"
	confData["logFile"] = "LOG_FILE"
	confData["logRetain"] = "LOG_RETAIN"
	data.ConfData = confData
	return command.CreateTemplateFile("config/app.go", command.ConfigTemplate, data)
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
	err = a.createHttpRouter(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 创建cmd/httpd/main.go文件
func (a Project) createHttpd(ctx *prowjob.Context) error {
	data := command.TemplateData{
		PackageName: "",
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
			{ImportName: moduleName + "/app/http/router"},
			{ImportName: moduleName + "/config"},
			{ImportName: frameworkModuleName + "/logger"},
			{Alias: "_", ImportName: frameworkModuleName + "/validator"},
			{ImportName: "github.com/spf13/cast"},
			{ImportName: "github.com/gin-gonic/gin"},
		},
	}
	return command.CreateTemplateFile("cmd/httpd/main.go", command.HttpdTemplate, data)
}

// 创建app/http/controller/demo.go文件
func (a Project) createHttpController(ctx *prowjob.Context) error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Index", Params: "ctx *gin.Context", Results: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller index demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Show", Params: "ctx *gin.Context", Results: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller show demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Update", Params: "ctx *gin.Context", Results: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller update demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Store", Params: "ctx *gin.Context", Results: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller store demo\"})"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Destroy", Params: "ctx *gin.Context", Results: "", FuncBody: "ctx.JSON(http.StatusOK, gin.H{\"message\": \"http controller destroy demo\"})"})

	data := command.TemplateData{
		PackageName: "controller",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/gin-gonic/gin"}, {ImportName: "net/http"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("app/http/controller/demo.go", command.CommonTemplate, data)
}

// 创建app/http/router/router.go文件
func (a Project) createHttpRouter(ctx *prowjob.Context) error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	var funcs []command.FuncTemplate

	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Register", Params: "eng *gin.Engine", Results: "", FuncBody: "apiGroup := eng.Group(\"/api\")\n\t{\n\trouter.ApiResource(apiGroup, \"/demo\", &controller.Demo{})\n\t}"})

	data := command.TemplateData{
		PackageName: "router",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/gin-gonic/gin"}, {ImportName: moduleName + "/app/http/controller"}, {ImportName: frameworkModuleName + "/http/restful/router"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("app/http/router/router.go", command.CommonTemplate, data)
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
	protocGenGoVersion := "v1.28.1"
	rs, err := execCommand("protoc-gen-go", "--version")
	if err != nil {
		return err
	}
	version := strings.Split(rs, " ")
	if len(version) < 2 {
		return errors.New("want protoc-gen-go version,but got " + rs)
	}

	if compareVersions(version[1], protocGenGoVersion) < 0 {
		return errors.New("it is recommended that protoc-gen-go be at least " + protocGenGoVersion)
	}
	return nil
}

func (a Project) protocGenGoGrpcVersion() error {
	protocGenGoGrpcVersion := "1.3.0"
	rs, err := execCommand("protoc-gen-go-grpc", "--version")
	if err != nil {
		return err
	}
	version := strings.Split(rs, " ")
	if len(version) < 2 {
		return errors.New("want protoc-gen-go-grpc version,but got " + rs)
	}

	if compareVersions(version[1], protocGenGoGrpcVersion) < 0 {
		return errors.New("it is recommended that protoc-gen-go-grpc be at least " + protocGenGoGrpcVersion)
	}
	return nil
}

func (a Project) buildProto() error {
	rs, err := execCommand("protoc", "-I=app/grpc/pb", "--go_out=app/grpc/pb", "--go-grpc_out=app/grpc/pb", "app/grpc/pb/demo.proto")
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
			{Alias: "grpcRouter", ImportName: moduleName + "/app/grpc/router"},
			{ImportName: moduleName + "/config"},
		},
	}

	return command.CreateTemplateFile("cmd/grpc/main.go", command.GrpcTemplate, data)
}

// 创建app/grpc/controller/demo.go文件
func (a Project) createGrpcController(ctx *prowjob.Context) error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Foo", Params: "context.Context, *demo.DemoReq", Results: "(*demo.DemoRes, error)", FuncBody: "return &demo.DemoRes{Name: \"Hello World\"}, nil"})

	data := command.TemplateData{
		PackageName: "controller",
		Imports:     []command.ImportTemplate{{ImportName: "context"}, {ImportName: "google.golang.org/grpc"}, {ImportName: moduleName + "/app/grpc/pb/demo"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName, Fields: []string{"Server *grpc.Server", "demo.UnimplementedDemoServer"}}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("app/grpc/controller/demo.go", command.CommonTemplate, data)
}

// 创建app/grpc/pb/demo.proto文件
func (a Project) createGrpcProto() error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Foo", Params: "DemoReq", Results: "DemoRes", FuncBody: ""})

	var msg []command.MessageTemplate
	msg = append(msg, command.MessageTemplate{MsgName: "DemoReq", Fields: []string{"string name = 1;"}})
	msg = append(msg, command.MessageTemplate{MsgName: "DemoRes", Fields: []string{"string name = 1;"}})

	data := command.TemplateData{
		PackageName: "/demo",
		Imports:     nil,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName}},
		Funcs:       funcs,
		ProtoPkg:    typeName,
		Services:    []command.ServiceTemplate{{Name: typeName, Funcs: funcs}},
		Messages:    msg,
	}
	return command.CreateTemplateFile("app/grpc/pb/demo.proto", command.ProtoTemplate, data)
}

// 创建app/grpc/router/router.go文件
func (a Project) createGrpcRouter() error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Register", Params: "s *grpc.Server", Results: "", FuncBody: "demo.RegisterDemoServer(s, &controller.Demo{Server: s})"})

	data := command.TemplateData{
		PackageName: "router",
		Imports:     []command.ImportTemplate{{ImportName: moduleName + "/app/grpc/controller"}, {ImportName: moduleName + "/app/grpc/pb/demo"}, {ImportName: "google.golang.org/grpc"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("app/grpc/router/router.go", command.CommonTemplate, data)
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
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "main", Params: "", Results: "", FuncBody: "jobEng := prowjob.New()\nprowjobreg.Register(jobEng)\nregister.Register(jobEng)\njobEng.Run()"})

	data := command.TemplateData{
		PackageName: "main",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}, {Alias: "prowjobreg", ImportName: frameworkModuleName + "/prowjob/register"}, {ImportName: moduleName + "/app/console/register"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("cmd/job/main.go", command.CommonTemplate, data)
}

// 创建app/console/command/demo.go文件
func (a Project) createCommand(ctx *prowjob.Context) error {
	typeName := "Demo"
	receiver := "d"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "GetCommand", Params: "", Results: "string", FuncBody: "return \"command:demo\""})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Usage", Params: "", Results: "string", FuncBody: "return `Usage of command:demo:\n\t  command:demo\n\t`"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Handle", Params: "ctx *prowjob.Context", Results: "", FuncBody: "fmt.Println(\"this is command demo\")"})

	data := command.TemplateData{
		PackageName: "command",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}, {ImportName: "fmt"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("app/console/command/demo.go", command.CommonTemplate, data)
}

// 创建app/console/command/register.go文件
func (a Project) createCommandRegister(ctx *prowjob.Context) error {
	typeName := ""
	receiver := ""
	receiverType := "" + typeName
	mn, err := module.GetName()
	if err != nil {
		return err
	}
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Register", Params: "eng *prowjob.CommandEngine", Results: "", FuncBody: "eng.Add(&command.Demo{})"})
	data := command.TemplateData{
		PackageName: "register",
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}, {ImportName: mn + "/app/console/command"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("app/console/register/register.go", command.CommonTemplate, data)
}

// 创建domain示例相关的文件
func (a Project) createDomainDemoFiles(ctx *prowjob.Context) error {
	err := a.createDemoAggRoot(ctx)
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
	typeName := "Root"
	receiver := "r"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "New", Params: "", Results: "*Root", FuncBody: "return &Root{}"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "GetTest", Params: "", Results: "string", FuncBody: "return \"arr root test\""})

	data := command.TemplateData{
		PackageName: "aggdemo",
		Imports:     nil,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile("domain/demo/aggroot.go", command.CommonTemplate, data)
}

// 创建domain/aggdemo/entity.go文件
func (a Project) createDemoAggEntity(ctx *prowjob.Context) error {
	typeName := "EntityA"
	receiver := "e"
	receiverType := "*" + typeName
	var funcs []command.FuncTemplate
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "NewA", Params: "", Results: "*EntityA", FuncBody: "return &EntityA{}"})
	funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "GetTest", Params: "", Results: "string", FuncBody: "return \"entity a test\""})

	data := command.TemplateData{
		PackageName: "aggdemo",
		Imports:     nil,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName}},
		Funcs:       funcs,
	}

	return command.CreateTemplateFile("domain/demo/entity.go", command.CommonTemplate, data)
}

// 创建domain/demoagg/vo.go文件
func (a Project) createDemoAggVo(ctx *prowjob.Context) error {
	typeName := "VoA"

	data := command.TemplateData{
		PackageName: "aggdemo",
		Imports:     nil,
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: typeName}},
		Funcs:       nil,
	}
	return command.CreateTemplateFile("domain/demoagg/vo.go", command.CommonTemplate, data)
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
	s, err = execCommand("go", "mod", "tidy")
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
