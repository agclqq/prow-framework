# validator
## intro
- the validator wraps translations and extends aliases on top of [validator](https://github.com/go-playground/validator/) V10
- translation
  - 17 languages are supported
- alias
  - Almost all the time the backend validates a field that is inconsistent with the front end, but the translator will not translate that field.  So the alias function is extended to accommodate the translation of that field.
  - Aliases are not translated.

## usage
### use with gin
Chinese translation and alias extension are enabled for gin framework in default mode. The default mode is completely non-invasive to gin.

```go
package Foo

import (
  "fmt"

  "github.com/agclqq/prow-framework/validator"

  "github.com/gin-gonic/gin"
)

type demoForm struct {
  Name string `form:"name" binding:"required" alias:"用户名"`
}

func Bar(ctx *gin.Context) {
	var df demoForm
	if err := ctx.ShouldBind(&df); err != nil {
		fmt.Println(validator.GetError(err))
		//or
		fmt.Println(validator.GetErrors(err))
	}
}
```
**Switch translation language**

Note: Only the default mode can be used to switch languages. Be sure to declare it before the validation method, and switching languages is global. It is recommended to put the switching language in the main package.
```go
validator.SwitchGinVldLang("zh_hant_tw")
...
ctx.ShouldBind(&foo)
```
### General use case
If you do not use gin, or do not want to use the verification methods provided by gin, then the following methods are recommended.
```go
package foo

import (
  "fmt"
  
  "github.com/agclqq/prow-framework/validator"
  validatorV10 "github.com/go-playground/validator/v10"
)

type demoForm struct {
  Name string `form:"name" binding:"required" alias:"用户名"`
}

func Bar() {
    tran,err:=validator.New(validatorV10.New(),validator.WithLocal("zh_hant_tw"),validator.WithAliasTag())
	err:=tran.Vld.Struct(demoForm{Name:""})
    if err != nil {
        fmt.Println(tran.GetError(err))
    }
}
```
