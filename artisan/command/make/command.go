package make

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/agclqq/prow/args"
	"github.com/agclqq/prow/artisan"
	"github.com/agclqq/prow/module"
	strings2 "github.com/agclqq/prow/strings"
)

const DEFAULT_COMMAND_DIR = "application/console/command/"

type Command struct {
}

func (a Command) GetCommand() string {
	return "make:command"
}

func (a Command) Usage() string {
	return `Usage of make:command:
  make:command commandName [path]
`
}
func (a Command) Handle(ctx *artisan.Context) {
	if len(ctx.Param) < 2 {
		fmt.Printf("%s \n %s", artisan.ERROR_PARAM_NUM, a.Usage())
		return
	}
	pparam := args.TidyParmaWithPrefix(ctx.Param)
	if _, ok := pparam["h"]; ok {
		fmt.Println(a.Usage())
		return
	}

	commandName := ctx.Param[1]
	path := ""
	if len(ctx.Param) == 3 {
		path = strings.TrimRight(ctx.Param[2], "/")
	}
	if !checkOverwrite(getCommandFullPath(commandName, path)) { //Abort if the file exists and the user does not allow it to be overwritten
		return
	}
	if err := createCommandFile(commandName, path); err != nil {
		fmt.Println(err)
	}
}
func getCommandFullPath(commandName, path string) string {
	return path + "/" + commandName + ".go"
}
func createCommandFile(commandName, path string) error {
	path = strings.TrimRight(path, "/")
	_, packageName := filepath.Split(path)
	receiver := strings2.ToLowFirst(commandName[0:1])
	receiverType := strings2.ToUpFirst(commandName)
	moduleName, err := module.GetModuleName()
	if err != nil {
		return err
	}
	data := TemplateData{
		PackageName:  packageName,
		ModuleName:   moduleName,
		Receiver:     receiver,
		ReceiverType: receiverType,
		CommandName:  commandName,
	}
	return CreateTemplateFile(getCommandFullPath(commandName, path), commandTemplate, data)
}
