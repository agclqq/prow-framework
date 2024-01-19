package list

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/prettytable"
)

const FORMAT_DEFAULT_LENGTH = 5

type Command struct {
}

func (a Command) GetCommand() string {
	return "list:command"
}
func (a Command) Usage() string {
	return `Usage of make:command:
  make:command commandName [path]
`
}
func (a Command) Handle(ctx *prowjob.Context) {
	commands := ctx.Engine.GetCommands()
	data := make([][]string, 0)
	data = append(data, []string{"Command", "Handler", "Desc"})
	for _, v := range commands {
		funcName := "func"
		if v.Ins != nil {
			funcName = getHandlerName(v.Ins)
		}

		data = append(data, []string{v.Command, funcName, strings.Replace(v.Desc, "\n", " ", -1)})
	}
	fmt.Println(prettytable.PlainText(data))
}

func getFileName(handler prowjob.Commander) string {
	t := reflect.TypeOf(handler)
	return t.PkgPath() + "/" + t.Name()
}
func getHandlerName(handler prowjob.Commander) string {
	// 获取实例的类型
	instanceType := reflect.TypeOf(handler)

	// 获取实例的包名和实例名称
	instancePkg := instanceType.PkgPath()
	instanceName := instanceType.Name()

	fmt.Println("Package:", instancePkg)
	fmt.Println("Instance Name:", instanceName)
	return instancePkg + "." + instanceName
}
