package make

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"

	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/args"
	"github.com/agclqq/prow-framework/file"
	"github.com/agclqq/prow-framework/module"
	"github.com/agclqq/prow-framework/prowjob/command"
	strings2 "github.com/agclqq/prow-framework/strings"
)

var defaultCommandDir = "app/console/command/"
var moduleName = ""

type Command struct {
}

func (a Command) GetCommand() string {
	return "make:command"
}

func (a Command) Usage() string {
	return `Usage of make:command:
  make:command commandName [path[ pathPrefix]] [-u] [-r]
    path is the path in the specified directory
    if the 'pathPrefix' is not given, the default is '` + defaultCommandDir + `'
	-u indicates the use case description
    -r [y/n] indicates whether to automatically register with the default registry file， the default is y
`
}
func (a Command) Handle(ctx *prowjob.Context) {
	if len(ctx.Param) < 1 {
		fmt.Println("error:" + command.NO_COMMAND_NAME)
		fmt.Println(a.Usage())
		return
	}
	prefixParam := args.TidyParmaWithPrefix(ctx.Param)
	if _, ok := prefixParam["h"]; ok {
		fmt.Println(a.Usage())
		return
	}

	mn, err := module.GetName()
	if err != nil {
		fmt.Println(err)
		return
	}
	moduleName = mn

	commandName := ctx.Param[0]
	path := ""
	if len(ctx.Param) == 2 {
		path = strings.TrimRight(ctx.Param[1], "/") + "/"
		if len(ctx.Param) == 3 {
			defaultCommandDir = ctx.Param[2]
		}
	}
	usage := ""
	if v, ok := prefixParam["u"]; ok {
		usage = fmt.Sprintf("Usage of %s:\n  %s", commandName, v)
	}

	fullCtlPath := GetFullPath(defaultCommandDir, path, commandName)
	if !CheckOverwrite(fullCtlPath) { //Abort if the file exists and the user does not allow it to be overwritten
		return
	}
	dir, _ := filepath.Split(fullCtlPath)
	_, packageName := filepath.Split(strings.TrimRight(dir, "/"))
	receiver := strings2.ToLowFirst(commandName[0:1])
	receiverType := strings2.ToUpFirst(commandName)

	if err = createCommandFile(packageName, receiver, receiverType, commandName, usage, fullCtlPath); err != nil {
		fmt.Println(err)
	}

	if v, ok := prefixParam["r"]; !ok || v == "" || v == "y" {
		err = registerCommand(defaultCommandRegister+"/register.go", strings.TrimRight(moduleName+"/"+dir, "/"), packageName, receiverType)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func createCommandFile(packageName, receiver, receiverType, commandName, usage, fullCtlPath string) error {
	data := command.TemplateData{
		PackageName: packageName,
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}},
		Consts:      nil,
		Vars:        nil,
		Types:       []command.TypeTemplate{{TypeName: receiverType}},
		Funcs: []command.FuncTemplate{
			{Receiver: receiver, ReceiverType: "*" + receiverType, FuncName: "GetCommand", Params: "", Results: "string", FuncBody: fmt.Sprintf("return \"command:%s\"", commandName)},
			{Receiver: receiver, ReceiverType: "*" + receiverType, FuncName: "Usage", Params: "", Results: "string", FuncBody: fmt.Sprintf("return `%s`", usage)},
			{Receiver: receiver, ReceiverType: "*" + receiverType, FuncName: "Handle", Params: "ctx *prowjob.Context", Results: "", FuncBody: ""},
		},
	}
	return command.CreateTemplateFile(fullCtlPath, command.CommonTemplate, data)
}

func registerCommand(commandRegisterFile, importVal, packageName, commandName string) error {
	// 读取文件，解析代码
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, commandRegisterFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// 找到相关的函数
	var importDecl *ast.GenDecl
	var registerFunc *ast.FuncDecl
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.GenDecl); ok && fd.Tok == token.IMPORT {
			importDecl = fd
		}
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == "Register" {
			registerFunc = fd
		}
	}

	// 添加新的导入路径和 ImportSpec
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s"`, importVal),
		},
	}
	importDecl.Specs = append(importDecl.Specs, newImport)
	sort.Slice(importDecl.Specs, func(i, j int) bool {
		return importDecl.Specs[i].(*ast.ImportSpec).Path.Value < importDecl.Specs[j].(*ast.ImportSpec).Path.Value
	})

	for _, v := range registerFunc.Body.List {
		if exprStmt, ok := v.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				if len(callExpr.Args) > 0 {
					if unaryExpr, ok := callExpr.Args[0].(*ast.UnaryExpr); ok {
						if compLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
							if selExpr, ok := compLit.Type.(*ast.SelectorExpr); ok {
								if selExpr.X.(*ast.Ident).Name == packageName && selExpr.Sel.Name == commandName {
									fmt.Printf("command.%s{} already exists \n", commandName)
									return nil
								}
							}
						}
					}
				}
			}
		}
	}
	// 在 Register 函数体中插入新行
	stmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "eng"},
				Sel: &ast.Ident{Name: "Add"},
			},
			Args: []ast.Expr{
				//&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`&%s`, commandObj)},
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: packageName},
							Sel: &ast.Ident{Name: commandName},
						},
					},
				},
			},
		},
	}
	registerFunc.Body.List = append(registerFunc.Body.List, stmt)

	// 格式化代码
	var buf bytes.Buffer
	if err = format.Node(&buf, fset, f); err != nil {
		return err
	}

	fmt.Println(buf.String())
	err = file.ReWriteString(commandRegisterFile, buf.String())
	if err != nil {
		return err
	}
	return nil
}
