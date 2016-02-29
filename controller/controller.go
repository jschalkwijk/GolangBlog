package controller

import (
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/jschalkwijk/GolangBlog/model/posts"
)

var view, _ = filepath.Abs("../jschalkwijk/GolangBlog/view")
var templates, _ = filepath.Abs("../jschalkwijk/GolangBlog/templates")

func RenderTemplate(w http.ResponseWriter,name string, data *posts.Data) {
	t, err := template.ParseFiles(templates+"/"+"header.html",templates+"/"+"nav.html",view + "/" + name + ".html",templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,data)
	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
