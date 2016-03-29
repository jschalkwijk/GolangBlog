package login

import (
	"net/http"
	"html/template"
)

var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

func RenderTemplate(w http.ResponseWriter,name string) {
	t, err := template.ParseFiles(templates+"/"+"header.html",view + "/" + name + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,name,nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}