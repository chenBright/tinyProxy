package web

import (
	"html/template"
	"log"
	"net/http"
)

// 渲染页面
func Render(writer http.ResponseWriter, templateName string, templateString string, context interface{}) {
	t := template.New(templateName)
	tpl, err := t.Parse(templateString)
	if err != nil {
		log.Println("render failed, template: %s", templateName)
		return
	}

	tpl.Execute(writer, context)
}

