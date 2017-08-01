package goblog

import (
	"html/template"
	"net/http"
)

var DefaultTemplates = []string{"base.html", "nav.html"}
var TemplateDirectory string = "templates/"

func init() {}

type View struct {
	tpl *template.Template
	layout string
	Title string
	Data map[string]interface{}
}

func (v *View) Execute(w http.ResponseWriter) {
	v.Data["Title"] = v.Title
	err := v.tpl.ExecuteTemplate(w, v.layout, v.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Execute(w)
}

func addPath(files ...string) []string {
	var f []string
	for _, file := range append(DefaultTemplates, files...) {
		f = append(f, TemplateDirectory + file)
	}
	return f
}

func NewView(layout string, title string, files ...string) (*View, error) {
	f := addPath(files...)
	tpl, err := template.ParseFiles(f...)
	if err != nil {
		return nil, err
	}

	return &View{
		layout: layout,
		tpl: tpl,
		Title: title,
		Data: make(map[string]interface{}),
	}, nil
}