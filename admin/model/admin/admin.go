package admin

import (
	"html/template"
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/posts"
	"github.com/jschalkwijk/GolangBlog/admin/model/users"
)

var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

type Data struct {
	Posts posts.Data
	Users users.Data
	Dashboard bool
}

func RenderTemplate(w http.ResponseWriter,name string, data *Data) {
	t, err := template.ParseFiles(templates+"/"+"header.html",templates+"/"+"nav.html",view + "/" + name + ".html",templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If you use the Parsefiles func to render the templates befirehand then you don't need to call them inside the index template.
	// or any other. This is just an example for me to keep this in mind.
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,data)
	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}