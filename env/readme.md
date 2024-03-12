# env
## intro
- env is a tool to manage your environment variables
- env files for multiple environments are currently supported

## useage
- put .env files in the root directory of your project

```go
package main

import (
	"fmt"

	"github.com/agclqq/prow-framework/env"
)

func main() {
	//std env, default is dot, globally unique instance
	fmt.Println(env.Get("TEST_KEY"))
	fmt.Println(env.Get("TEST_KEY2", "default value"))

	// multiple instances
	em, err := env.New(env.Dot)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(em.Get("TEST_KEY"))

	//files are determined by specifying the environment
	//The default environment variable name is GO_ENV
	em, err := env.New(env.Dot, env.WithOsEnv())
	...

	//you can also specify environment variable names
	em, err := env.New(env.Dot, env.WithOsEnv(), env.WithEnvName("GO_ENV"))
	...

	//you can also specify the file name
	em, err := env.New(env.Dot, env.WithFile(".env"))
	...

	//when you use multiple instances,you can override these envs
	em.Set("TEST_KEY", "new value")
	em.SetAll(map[string]string{"TEST_KEY": "new value", "TEST_KEY2": "new value2"})
}
```