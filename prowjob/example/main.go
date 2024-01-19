package main

import (
	"github.com/agclqq/prowjob"

	"github.com/agclqq/prow-framework/prowjob/register"
)

func main() {
	pj := prowjob.New()
	register.Register(pj)
	pj.Run("init:project")
	//pj.Run("make:command", "test2")
	//pj.Run("make:controller", "test")
	//pj.Run("make:commandRegister")
}
