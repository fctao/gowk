package view

import (
	"html/template"
	"net/http"
)

type ModelAndView struct {
	ViewName  string
	TplNames  string
	ModelName string
	Model     interface{}
}

func (mav *ModelAndView) Execute(writer http.ResponseWriter) error {
	tpl, err := template.New(mav.ViewName).ParseFiles(mav.TplNames)
	if nil != err {
		return err
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if mav.ModelName == "" {
		tpl.Execute(writer, mav.Model)
	} else {
		tpl.ExecuteTemplate(writer, mav.ModelName, mav.Model)
	}
	return nil
}
