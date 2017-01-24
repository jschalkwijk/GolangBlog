package blog

import (
	_"github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/posts"
)

// here we define the absolute path to the view folder it takes the go root until the github folder.
var view = "GolangBlog/view"
var templates = "GolangBlog/templates"

/*
  The function template.ParseFiles will read the contents of "".html and return a *template.Template.
  The method t.Execute executes the template, writing the generated HTML to the http.ResponseWriter.
  The .Title and .Body dotted identifiers inside the template refer to p.Title and p.Body.
*/


func View(w http.ResponseWriter,name string, p *posts.Data) {
	t, err := template.ParseFiles(templates+"/"+"header.html",templates+"/"+"nav.html",view + "/" + name + ".html",templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,p)
	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
