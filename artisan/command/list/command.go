package list

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/agclqq/prow/artisan"
)

const FORMAT_DEFAULT_LENGTH = 5

type Command struct {
}

func (a Command) GetCommand() string {
	return "list:command"
}
func (a Command) Handle(ctx *artisan.Context) {
	commands := ctx.Engine.GetCommands()
	format := getCommandLenForFormat(commands)
	res := make([]string, 0)
	for _, v := range commands {
		res = append(res, fmt.Sprintf(format, v.Name, getFileName(v.Handler), v.Desc))
		getFileName(v.Handler)
	}
	sort.Strings(res)
	fmt.Println(strings.Join(res, ""))
}

func getCommandLenForFormat(commands map[string]artisan.Command) string {
	nameList := make([]string, 10)
	pathList := make([]string, 10)
	descList := make([]string, 10)
	for _, v := range commands {
		nameList = append(nameList, v.Name)
		pathList = append(pathList, getFileName(v.Handler))
		descList = append(descList, v.Desc)
	}
	nameLen := getMaxLenFromList(nameList)
	pathLen := getMaxLenFromList(pathList)
	descLen := getMaxLenFromList(descList)
	s := fmt.Sprintf("%%-%ds|%%-%ds|%%-%ds\n", nameLen, pathLen, descLen)
	return s
}

func getMaxLenFromList(list []string) int {
	maxLen := FORMAT_DEFAULT_LENGTH
	for _, v := range list {
		if len(v) > maxLen {
			maxLen = len(v)
		}
	}
	return maxLen
}
func getFileName(handler artisan.Commander) string {
	t := reflect.TypeOf(handler)
	return t.PkgPath() + "/" + t.Name()
}
