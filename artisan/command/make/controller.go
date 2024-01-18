package make

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/agclqq/prow/artisan"
	file2 "github.com/agclqq/prow/file"
	"github.com/agclqq/prow/module"
	strings2 "github.com/agclqq/prow/strings"

	"github.com/agclqq/prow/args"
)

const DEFAULT_CONTROLLER_DIR = "application/http/controller"
const ROUTER_DIR = "application/http/router/"

type Controller struct {
}

func (a Controller) GetCommand() string {
	return "make:controller"
}
func (a Controller) Usage() string {
	return `Usage of make:controller:
  make:controller controllerName [path] [-r y/n]
    If the 'path' is not given, the default is '` + DEFAULT_CONTROLLER_DIR + `'
    if the '-r' is not given,the default is y
`
}
func (a Controller) Handle(ctx *artisan.Context) {
	fmt.Println(ctx.Param)
	param := tidyParmaNoPrefix(ctx.Param)
	pparam := args.TidyParmaWithPrefix(ctx.Param)
	if _, ok := pparam["h"]; ok {
		fmt.Println(a.Usage())
		return
	}
	//make:controller ctlName [path] [-r n]
	if len(param) <= 1 {
		fmt.Println("error:" + artisan.NO_CONTROLLER_NAME)
		return
	}
	ctlName := param[1]
	ctlPath := ""
	if len(param) >= 3 {
		ctlPath = param[2]
	}
	ctlType := strings.ToLower(pparam["r"])
	if ctlType == "" {
		ctlType = "y"
	}
	if !checkOverwrite(getControllerFullPath(ctlName, ctlPath)) { //Abort if the file exists and the user does not allow it to be overwritten
		return
	}

	if err := createControllerFile(ctlName, ctlPath, ctlType); err != nil {
		fmt.Printf("error of %s:%s", param, err)
	}
	//createRouter(ROUTER_DIR + "router.go")
}
func getControllerFullPath(ctlName, ctlPath string) string {
	return ctlPath + "/" + strings.ToLower(ctlName) + ".go"
}

func checkOverwrite(file string) bool {
	if file2.Exist(file) {
		fmt.Printf(artisan.CONTROLLER_EXIST+" \n", file)
		fmt.Println("whether to overwrite this file?[y/n]")
		goon := ""
		if _, err := fmt.Scanf("%s", &goon); err != nil {
			fmt.Println(err)
			return false
		}
		if goon != "y" {
			return false
		}
	}
	return true
}

// 整理参数列表，只返回无前导符(-，--)的参数
func tidyParmaNoPrefix(param []string) []string {
	newParam := make([]string, 0)
	nextSkip := false
	for _, v := range param {
		if (strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--")) && !strings.Contains(v, "=") {
			nextSkip = true
			continue
		}
		if nextSkip {
			continue
		}
		newParam = append(newParam, v)
	}
	return newParam
}

func createRouter(fileName string) error {
	//fileName := ROUTER_DIR + "router.go"
	err := file2.Touch(fileName)
	if err != nil {
		return err
	}

	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, fileName, nil, parser.Trace)
	if err != nil {
		panic(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		//var l []string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		case *ast.FuncDecl:
			s = x.Name.Name
		case *ast.FuncType:
			fmt.Println(x.Func)
		case *ast.FieldList:
			fmt.Println(x.List)
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fSet.Position(n.Pos()), s)
		}
		return true
	})
	fmt.Println("over")
	return nil
}
func createControllerFile(ctlName, ctlPath, ctlType string) error {
	ctlPath = strings.TrimRight(ctlPath, "/")
	_, packageName := filepath.Split(ctlPath)
	moduleName, err := module.GetModuleName()
	if err != nil {
		return err
	}
	receiver := strings2.ToLowFirst(ctlName[0:1])
	receiverType := strings2.ToUpFirst(ctlName)
	isResource := false
	if ctlType == "y" {
		isResource = true
	}
	data := TemplateData{
		PackageName:  packageName,
		ModuleName:   moduleName,
		Receiver:     receiver,
		ReceiverType: receiverType,
		IsResource:   isResource,
	}
	return CreateTemplateFile(getControllerFullPath(ctlName, ctlPath), controllerTemplate, data)
}
