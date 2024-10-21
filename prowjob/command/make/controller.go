package make

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/agclqq/prowjob"

	strings2 "github.com/agclqq/prow-framework/strings"

	"github.com/agclqq/prow-framework/args"
	"github.com/agclqq/prow-framework/prowjob/command"
)

var defaultControllerDir = "app/http/controller/"

const ROUTER_DIR = "app/http/router/"

type Controller struct {
}

func (a Controller) GetCommand() string {
	return "make:controller"
}

func (a Controller) Usage() string {
	return `Usage of make:controller:
  make:controller controllerName [path[ pathPrefix]] [-r y/n]
    the path refers to the part after the prefix
    if the 'pathPrefix' is not given, the default is '` + defaultControllerDir + `'
    -r indicates whether it is restful api, the default is y
`
}

func (a Controller) Handle(ctx *prowjob.Context) {
	if len(ctx.Param) < 1 {
		fmt.Println("error:" + command.NO_CONTROLLER_NAME)
		return
	}
	param := args.TidyParmaNoPrefix(ctx.Param)
	PrefixParma := args.TidyParmaWithPrefix(ctx.Param)
	if _, ok := PrefixParma["h"]; ok {
		fmt.Println(a.Usage())
		return
	}
	//make:controller ctlName [path] [-r n]
	ctlName := param[0]
	ctlPath := ""
	if len(param) >= 2 {
		ctlPath = strings.TrimRight(param[1], "/") + "/"
		if len(param) >= 3 {
			defaultControllerDir = param[2]
		}
	}
	ctlType := strings.ToLower(PrefixParma["r"])
	isResource := true
	if ctlType == "n" {
		isResource = false
	}
	fullCtlPath := GetFullPath(defaultControllerDir, ctlPath, ctlName)
	if !CheckOverwrite(fullCtlPath) { //Abort if the file exists and the user does not allow it to be overwritten
		return
	}

	dir, _ := filepath.Split(fullCtlPath)
	_, packageName := filepath.Split(strings.TrimRight(dir, "/"))
	receiver := strings2.ToLowFirst(ctlName[0:1])
	receiverType := strings2.ToUpFirst(ctlName)
	if err := createControllerFile(packageName, receiver, receiverType, fullCtlPath, isResource); err != nil {
		fmt.Printf("error of %s:%s", param, err)
	}
}

// packageName, receiver, receiverType, commandName, usage, fullCtlPath
func createControllerFile(packageName, receiver, receiverType, ctlPath string, isResource bool) error {
	var funcs []command.FuncTemplate
	if isResource {
		funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Index", Params: "ctx *gin.Context", ResultType: "", FuncBody: ""})
		funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Show", Params: "ctx *gin.Context", ResultType: "", FuncBody: ""})
		funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Store", Params: "ctx *gin.Context", ResultType: "", FuncBody: ""})
		funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Update", Params: "ctx *gin.Context", ResultType: "", FuncBody: ""})
		funcs = append(funcs, command.FuncTemplate{Receiver: receiver, ReceiverType: receiverType, FuncName: "Destroy", Params: "ctx *gin.Context", ResultType: "", FuncBody: ""})
	}
	data := command.TemplateData{
		PackageName: packageName,
		Imports:     []command.ImportTemplate{{ImportName: "github.com/gin-gonic/gin"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: receiverType}},
		Funcs:       funcs,
	}
	return command.CreateTemplateFile(ctlPath, command.CommonTemplate, data)
}
