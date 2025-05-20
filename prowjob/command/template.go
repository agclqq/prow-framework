package command

import (
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const pkg = `package {{.PackageName}}`
const imports = `{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	{{if .Alias}}{{.Alias}}{{end}} "{{.ImportName}}"{{end}}
)
{{end}}`

const consts = `{{if gt (len .Consts) 0}}
const ({{range .Consts}}
	{{.}}
{{end}}
){{end}}`

const vars = `{{if gt (len .Vars) 0}}
var ({{range .Vars}}
	{{.}}{{end}}
){{end}}`

const types = `{{if gt (len .Types) 0}}{{range .Types}}
type {{.Name}} struct { {{range .Fields}}
	{{.}}{{end}}
}{{end}}
{{end}}`
const interfaces = `{{if gt (len .Interfaces) 0}}{{range .Interfaces}}
type {{.Name}} interface {
	{{range .Methods}}{{.}}
	{{end}}
}{{end}}
{{end}}`

const funcs = `{{if gt (len .Funcs) 0}}{{range .Funcs}}
func {{if .Receiver}}({{.Receiver}} {{.ReceiverType}}) {{end}}{{.FuncName}}({{.Params}}) {{.ResultType}} {
	{{.FuncBody}}
}{{end}}
{{end}}`

const makefileVars = `{{range .Vars}}{{.}}
{{end}}`
const makefileRules = `{{if gt (len .MakefileRules) 0}}{{range .MakefileRules}}
{{.Target}}: {{.Dependencies}}{{if gt (len .Commands) 0}}
	{{range .Commands}}{{.}}
	{{end}}
{{end}}{{end}}{{end}}
`
const TextLines = `{{range $text := .TextLines}}{{.}}
{{end}}`
const CommonTemplate = pkg + imports + consts + vars + types + interfaces + funcs
const TmplTemplate = `
{{.ViewTemplateDefine}}
this is view template.
{{.ViewTemplateEnd}}
`

const MakefileTemplate = makefileVars + makefileRules

const EnvTemplate = `{{range $env := .Envs}}{{if eq $env.Type "comment"}}{{$env.Val}}{{else}}{{$env.Key}}={{$env.Val}}{{end}}
{{end}}`
const ConfigSingleMapTemplate = `
package config
` + imports + consts + vars + `
var {{.ConfData.ConfName}} = map[string]string{
{{range $key, $value := .ConfData.Vars}}	"{{$key}}":	autoenv.Get("{{$value}}"),
{{end}}
}

func Get{{.ConfData.ConfName | title}}(key string) string {
	return app[key]
}	
`
const ConfigDoubleMapTemplate = `
package config
` + imports + consts + vars + `
var {{.ConfData.ConfName}} = map[string]map[string]string{
{{range $key, $value := .ConfData.VarsM}}	"{{$key}}":	{
        {{range $k, $v:=$value}}"{{$k}}":autoenv.Get("{{$v}}"),
		{{end}}    },
{{end}}
}
func GetAll{{.ConfData.ConfName | title}}() map[string]map[string]string {
	return {{.ConfData.ConfName}}
}
func Get{{.ConfData.ConfName | title}}(key string) map[string]string {
	return {{.ConfData.ConfName}}[key]
}
`
const HttpdTemplate = pkg + `
` + imports + `

func main() {
	wg := &sync.WaitGroup{}

	//wg.Add(1)
	//go pprofServer(wg)

	wg.Add(1)
	go httpServer(wg)

	wg.Wait()
}

func pprofServer(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("start pprofServer at: %s\n", "6060")
	server := &http.Server{
		Addr:    ":6060",
		Handler: nil,
	}
	go func() {
		fmt.Printf("start pprofServer at: %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println(err)
				return
			}
			_ = fmt.Errorf("start pprofServer is error: %s\n", err)
		}
	}()
}

func httpServer(wg *sync.WaitGroup) {
	defer wg.Done()
	if config.GetApp("appEnv") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	//注册事件
	boot.StartEvent()

	eng := gin.New()
	eng.RedirectTrailingSlash = false
	eng.Use(logger.WithConfig(logger.AccessLogConfig(eng,config.GetApp("logFile"),cast.ToInt(config.GetApp("logRetain")))), gin.Recovery())
	router.Register(eng)
	eng.StaticFS("/resource", gin.Dir("./resource", false))
	eng.LoadHTMLGlob("resource/views/**/*")
	server := &http.Server{
		Addr:              ":" + config.GetApp("httpServerPort"),
		Handler:           eng,
		IdleTimeout:       75 * time.Second,
	}
	go func() {
		//if err := server.ListenAndServeTLS("resource/cert.pem", "resource/cert.key"); err != nil {
		//	fmt.Printf("start https server is error: %s\n", err)
		//}
		//(&provider.Event{}).Run()
		fmt.Printf("start http server at: %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println(err)
				return
			}
			fmt.Errorf("start http server is error: %s\n", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sign := <-ch
	fmt.Println("got a sign:", sign)
	now := time.Now()
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(cxt)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	// 看看实际退出所耗费的时间
	fmt.Println("http server is exited,cost:", time.Since(now).Milliseconds(), "ms")
}
`
const GrpcTemplate = `package main
` + imports + `

func main() {
	wg := &sync.WaitGroup{}

	//wg.Add(1)
	//go pprofServer(wg)

	wg.Add(1)
	go grpcServer(wg)

	wg.Wait()
}
func pprofServer(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("start pprofServer at: %s\n", "6060")
	server := &http.Server{
		Addr:    ":6060",
		Handler: nil,
	}
	go func() {
		fmt.Printf("start pprofServer at: %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println(err)
				return
			}
			fmt.Printf("start pprofServer is error: %s\n", err)
		}
	}()
}

func grpcServer(wg *sync.WaitGroup) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":" + config.GetApp("grpc_server_port"))
	if err != nil {
		_ = fmt.Errorf("failed to listen: %v", err)
	}
	
	kp := keepalive.ServerParameters{
		Time:    20 * time.Second,
		Timeout: 5 * time.Second,
	}
	kep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(kp),
		grpc.KeepaliveEnforcementPolicy(kep),
	)

	grpcRouter.Register(s)

	go func() {
		fmt.Printf("start grpc server at: %s\n", lis.Addr().String())
		if err = s.Serve(lis); err != nil {
			fmt.Printf("start grpc server is error: %v", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sign := <-ch
	fmt.Println("got a sign:", sign)
	now := time.Now()
	s.GracefulStop()
	// 看看实际退出所耗费的时间
	fmt.Println("grpc server is exited,cost:", time.Since(now).Milliseconds(), "ms")
}
`

const ProtoTemplate = `syntax = "proto3";
option go_package = "{{.PackageName}}";
package {{.ProtoPkg}};
{{range .Services}}
service {{.Name}} {
  {{range .Funcs}}rpc {{.FuncName}}({{.Params}}) returns({{.ResultType}});{{end}}
}{{end}}
{{range .Messages}}
message {{.MsgName}} {
  {{range .Fields}}{{.}}{{end}}
}{{end}}
`

type FuncTemplate struct {
	Receiver     string
	ReceiverType string
	FuncName     string
	Params       string
	ResultType   string
	FuncBody     string
}
type InterTemplate struct {
	Name    string
	Methods []string
}
type TypeTemplate struct {
	Name   string
	Fields []string
}
type ImportTemplate struct {
	Alias      string
	ImportName string
}
type ServiceTemplate struct {
	Name  string
	Funcs []FuncTemplate
}
type MessageTemplate struct {
	MsgName string
	Fields  []string
}
type Kv map[string]string
type KvM map[string]map[string]string
type EventTemplate struct {
	Key  string
	Val  string
	Type string
}
type TextLineData struct {
	TextLines []string
}
type ConfTemplate struct {
	ConfName string
	Vars     Kv
	VarsM    KvM
}
type MakefileData struct {
	Vars          []string
	MakefileRules []MakefileRule
}
type MakefileRule struct {
	Target       string
	Dependencies string
	Commands     []string
}
type TemplateData struct {
	PackageName        string
	Imports            []ImportTemplate
	Consts             []string
	Vars               []string
	Interfaces         []InterTemplate
	Types              []TypeTemplate
	Funcs              []FuncTemplate
	ProtoPkg           string
	Services           []ServiceTemplate
	Messages           []MessageTemplate
	CommandName        string
	IsResource         bool
	Envs               []EventTemplate
	ConfData           ConfTemplate
	MakefileRules      []MakefileRule
	ViewTemplateDefine string
	ViewTemplateEnd    string
}

func CreateTemplateFile(filePath string, tpl string, data any) error {
	filePath = filepath.Clean(filePath)
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 解析模板
	tmpl, err := template.New("GoTemplate").Funcs(template.FuncMap{
		"title": func(s string) string {
			s = cases.Title(language.English).String(s)
			return s
		},
	}).Parse(tpl)
	if err != nil {
		return err
	}
	// 执行模板，将结果写入文件
	err = tmpl.Execute(file, data)
	return err
}
