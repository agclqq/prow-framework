package artisan

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type CommandEngine struct {
	Commands map[string]Command
}
type Param struct {
	Key   string
	Value string
}

type Params []string
type Context struct {
	context.Context
	Param     Params
	Engine    *CommandEngine
	TidyParma map[string]string
}

type Command struct {
	Command     string
	HandlerFunc CommandFunc
	Desc        string
}

type CommandFunc func(ctx *Context)

func New() *CommandEngine {
	return &CommandEngine{Commands: make(map[string]Command)}
}
func (e *CommandEngine) Add(commander Commander, args ...string) {
	c := Command{
		Command:     commander.GetCommand(),
		HandlerFunc: commander.Handle,
		Desc:        commander.Usage(),
	}
	if len(args) > 0 {
		c.Desc = args[0]
	}
	e.Commands[c.Command] = c
}
func (e *CommandEngine) AddFunc(command string, f CommandFunc, args ...string) {
	c := Command{
		Command:     command,
		HandlerFunc: f,
	}
	if len(args) > 0 {
		c.Desc = args[0]
	}
	e.Commands[c.Command] = c

}
func (e *CommandEngine) Run(param ...string) {
	ctx := context.Background()
	e.RunWithCtx(ctx, param...)
}
func (e *CommandEngine) RunWithCtx(ctx context.Context, param ...string) {
	if !argsCheck(param...) {
		return
	}
	if len(param) == 0 {
		param = os.Args[1:]
	}
	e.Invoke(ctx, param)
}
func (e *CommandEngine) Invoke(ctx context.Context, param []string) {
	cmd, ok := e.Commands[param[0]]
	if !ok {
		fmt.Printf("error:"+NO_COMMAND, param[0])
		return
	}
	if cmd.Command != param[0] {
		fmt.Printf("error:"+INCONFORMITY, param[0], cmd.Command)
		return
	}
	newParam := make([]string, len(param)-1)
	if len(param) > 1 {
		newParam = param[1:]
	}
	c := &Context{
		Context:   ctx,
		Param:     newParam,
		TidyParma: TidyParam(newParam),
		Engine:    e,
	}
	//执行业务逻辑
	cmd.HandlerFunc(c)
}

func argsCheck(param ...string) bool {
	if len(os.Args) <= 1 && len(param) == 0 {
		fmt.Printf("error: %s", NO_PARAM)
		return false
	}
	return true
}

func (e *CommandEngine) GetCommands() map[string]Command {
	return e.Commands
}

// TidyParam  return the parameters with leading symbols (-, --) as well as those without
// e.g：-arg1 1 arg2 --arg3=3 -arg4=4 arg7 7 --arg5=5 arg6=6
// return：map[string]string{"arg1":"1","arg2":"","arg3":"3","arg4":"4","arg5":"5","arg6":"6","arg7":"7"}
func TidyParam(params []string) map[string]string {
	arrParams := strings.Split(strings.Join(strings.Split(strings.Join(params, " "), "="), " "), " ")
	arrParamsLen := len(arrParams)
	mapParams := make(map[string]string)
	for i := 0; i < arrParamsLen; i++ {
		for j := 0; j < 2; j++ {
			if arrParams[i] = strings.TrimLeft(arrParams[i], "-"); arrParams[i] == "" {
				continue
			}
		}
		mapParams[arrParams[i]] = ""
		if i+1 < arrParamsLen && (!strings.HasPrefix(arrParams[i+1], "-") && !strings.HasPrefix(arrParams[i+1], "--")) {
			mapParams[arrParams[i]] = arrParams[i+1]
			i++
		}
	}
	return mapParams
}
