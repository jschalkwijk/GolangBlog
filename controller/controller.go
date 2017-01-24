package controller

import (
	"html/template"
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/config"
)

func View(w http.ResponseWriter, name string, data interface{}){
	t, err := template.ParseFiles(config.Templates+"/"+"header.html",config.Templates+"/"+"nav.html",config.View + "/" + name + ".html",config.Templates+"/"+"footer.html")
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
