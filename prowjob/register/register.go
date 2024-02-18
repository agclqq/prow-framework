package register

import (
	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/prowjob/command/list"
	"github.com/agclqq/prow-framework/prowjob/command/make"
	"github.com/agclqq/prow-framework/prowjob/command/project"
)

func Register(eng *prowjob.CommandEngine) {
	eng.Add(project.Project{})
	eng.Add(make.Controller{}, "create controller, usage:artisan make:controller controllerName path [notresource]")
	eng.Add(make.Model{}, "create model, usage:artisan make:model modelName path")
	eng.Add(make.Command{})
	eng.Add(list.Command{})
	eng.Add(make.CommandRegister{})
}
