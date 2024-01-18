package register

import (
	"github.com/agclqq/prow/artisan"
	"github.com/agclqq/prow/artisan/command/list"
	"github.com/agclqq/prow/artisan/command/make"
)

func Register(eng *artisan.CommandEngine) {
	eng.Add("make:controller", make.Controller{}, "create controller, usage:artisan make:controller controllerName path [notresource]")
	//eng.Add("make:model", make.Model{}, "create model, usage:artisan make:model modelName path")
	eng.Add("make:command", make.Command{})
	eng.Add("list:command", list.Command{})
}
