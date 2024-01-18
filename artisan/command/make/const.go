package make

import (
	"os"
	"path/filepath"
	"text/template"
)

const commandTemplate = `
package {{.PackageName}}

import (
	"github.com/agclqq/prow-framework/artisan"
)

type {{.ReceiverType}} struct {
}

func ({{.Receiver}} {{.ReceiverType}}) GetCommand() string {
	return "command:{{.CommandName}}"
}

func ({{.Receiver}} {{.ReceiverType}}) Handle(ctx *artisan.Context) {

}
`
const controllerTemplate = `
package {{.PackageName}}

import (
	"github.com/gin-gonic/gin"
)

type {{.ReceiverType}} struct {
}
{{if .IsResource}}
func ({{.Receiver}} {{.ReceiverType}}) Index(ctx *gin.Context) {
	
}

func ({{.Receiver}} {{.ReceiverType}}) Store(ctx *gin.Context) {
	
}

func ({{.Receiver}} {{.ReceiverType}}) Show(ctx *gin.Context) {
	
}

func ({{.Receiver}} {{.ReceiverType}}) Update(ctx *gin.Context) {
	
}

func ({{.Receiver}} {{.ReceiverType}}) Destroy(ctx *gin.Context) {
	
}
{{end}}
`

type TemplateData struct {
	PackageName  string
	ModuleName   string
	Receiver     string
	ReceiverType string
	CommandName  string
	IsResource   bool
}

func CreateCommandTemplateFile(filePath string, data TemplateData) error {
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 解析模板
	tmpl, err := template.New("GoTemplate").Parse(commandTemplate)
	if err != nil {
		return err
	}
	// 执行模板，将结果写入文件
	err = tmpl.Execute(file, data)
	return err
}

func CreateTemplateFile(filePath string, tpl string, data TemplateData) error {
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 解析模板
	tmpl, err := template.New("GoTemplate").Parse(tpl)
	if err != nil {
		return err
	}
	// 执行模板，将结果写入文件
	err = tmpl.Execute(file, data)
	return err
}
