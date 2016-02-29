package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jschalkwijk/GolangBlog/model/home"
	"github.com/jschalkwijk/GolangBlog/controller/posts"
	"github.com/jschalkwijk/GolangBlog/controller/categories"
	"fmt"
	"github.com/jschalkwijk/GolangBlog/admin/model/admin"
//	"github.com/jschalkwijk/GolangBlog/admin/controller/posts"
//	"github.com/jschalkwijk/GolangBlog/admin/controller/categories"
)

var static = "/GolangBlog/static/"
var css = "/GolangBlog/static/css/"

func main() {
// With this funtion I can check if my filepath is working for serving static files such as CSS or Templates etc
// IMPORTANT:I failed to add stat ic files because Go will use the current Directory you are in as the App's ROOT.
// If I run it from GolangBlog, the root is /Users/jorn/Documents/Golang/src/github.com/jschalkwijk/GolangBlog
// If I run it from jschalkwijk
	// 	_, err := os.Stat(filepath.Join(".", "GolangBlog/static/css", "style.css"))
//	checkErr(err)
	fmt.Println("Starting GolangBlog..")
	r := mux.NewRouter()
	serveStatic(r,static)
	//r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("."+css))))

	r.HandleFunc("/", home.DashboardHandler)
	// Posts
	r.HandleFunc("/posts", posts.Index)
		p := r.PathPrefix("/posts").Subrouter()
		p.HandleFunc("/{id:[0-9]+}/{title}", posts.Single)
		p.HandleFunc("/add-post", posts.New)
		p.HandleFunc("/edit/{id:[0-9]+}/{title}", posts.Edit)
		p.HandleFunc("/save/{id:[0-9]+}/{title}", posts.Save)
		p.HandleFunc("/add", posts.Add)
	// Categories
	r.HandleFunc("/categories", categories.Index)
		c := r.PathPrefix("/categories").Subrouter()
		c.HandleFunc("/{id:[0-9]+}/{title}", categories.Single)
		c.HandleFunc("/add-category", categories.New)
		c.HandleFunc("/edit/{id:[0-9]+}/{title}", categories.Edit)
		c.HandleFunc("/save/{id:[0-9]+}/{title}", categories.Save)
		c.HandleFunc("/add", categories.Add)
	r.HandleFunc("/admin", admin.DashboardHandler)

	http.Handle("/", r)
	fmt.Println("Succes!")
	fmt.Println("GolangBlog running on port 8080. Don't forget to run MAMP or SQL server.")
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// load all static diretories: Source: http://www.shakedos.com/2014/Feb/08/serving-static-files-with-go.html
func serveStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"css":           staticDirectory + "/css/",
		"images":           staticDirectory + "/images/",
		"scripts":          staticDirectory + "/scripts/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		//fmt.Println(pathPrefix)
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir("."+pathValue))))
	}
}