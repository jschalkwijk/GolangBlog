package home

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var view, _ = filepath.Abs("../jschalkwijk/GolangBlog/view")
var templates, _ = filepath.Abs("../jschalkwijk/GolangBlog/templates")

type Collection struct {

}

func DashboardHandler(w http.ResponseWriter, r *http.Request){
		//params := splitURL(r,"")
		collection := new(Collection)
		p := collection
		renderTemplate(w,"index", p)
}

func renderTemplate(w http.ResponseWriter,name string, data *Collection) {
	t, err := template.ParseFiles(templates+"/"+"header.html",templates+"/"+"nav.html",view + "/" + name + ".html",templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If you use the Parsefiles func to render the templates befirehand then you don't need to call them inside the index template.
	// or any other. This is just an example for me to keep this in mind.
//	t.ExecuteTemplate(w,"header",nil)
//	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,data)
//	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}