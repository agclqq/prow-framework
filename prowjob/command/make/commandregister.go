package make

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/args"
	"github.com/agclqq/prow-framework/prowjob/command"
)

var defaultCommandRegister = "app/console/register"

type CommandRegister struct {
}

func (c CommandRegister) GetCommand() string {
	return "make:commandRegister"
}

func (c CommandRegister) Usage() string {
	return `Usage of make:commandRegister:
  make:commandRegister [path[ pathPrefix]]
    path is the path in the specified directory
    if the 'pathPrefix' is not given, the default is '` + defaultCommandRegister + `'
`
}

func (c CommandRegister) Handle(ctx *prowjob.Context) {
	prefixParam := args.TidyParmaWithPrefix(ctx.Param)
	if _, ok := prefixParam["h"]; ok {
		fmt.Println(c.Usage())
		return
	}

	commandName := "register"
	path := ""
	if len(ctx.Param) == 1 {
		path = strings.TrimRight(ctx.Param[0], "/") + "/"
		if len(ctx.Param) == 2 {
			defaultCommandDir = ctx.Param[1]
		}
	}

	fullPath := GetFullPath(defaultCommandRegister, path, commandName)
	if !CheckOverwrite(fullPath) { //Abort if the file exists and the user does not allow it to be overwritten
		return
	}
	if err := createCommandRegisterFile(commandName, fullPath); err != nil {
		fmt.Println(err)
	}
}

func createCommandRegisterFile(commandName, commandPath string) error {
	dir, _ := filepath.Split(commandPath)
	_, packageName := filepath.Split(strings.TrimRight(dir, "/"))
	//receiver := strings2.ToLowFirst(commandName[0:1])
	//receiverType := strings2.ToUpFirst(commandName)
	//_, err := module.GetModuleName()
	//if err != nil {
	//	return err
	//}
	data := command.TemplateData{
		PackageName: packageName,
		Imports:     []command.ImportTemplate{{ImportName: "github.com/agclqq/prowjob"}},
		Consts:      nil,
		Vars:        nil,
		Types:       nil,
		Funcs: []command.FuncTemplate{
			{FuncName: "Register", Params: "eng *prowjob.CommandEngine", FuncBody: ""},
		},
	}
	return command.CreateTemplateFile(commandPath, command.CommonTemplate, data)
}
