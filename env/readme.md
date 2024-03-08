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
	fmt.Println(env.Get("TEST_KEY"))
	fmt.Println(env.Get("TEST_KEY2","default value"))
	
	// or
	em, err := env.New(env.Dot)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(em.Get("TEST_KEY"))
}
```