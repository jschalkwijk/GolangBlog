package users

import (
	"html/template"
	"net/http"
)

var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

type User struct {
	User_ID int
	Username string
	FirstName string
	LastName string
	Keywords string
	Approved int
	Author string
	Date string
	Category_ID int
	Category string
	Trashed int
}

type Data struct {
	Users      []User
	Deleted	bool
}

func RenderTemplate(w http.ResponseWriter,name string, p *Data) {
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
