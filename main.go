package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jschalkwijk/GolangBlog/model/home"
	"github.com/jschalkwijk/GolangBlog/controller/posts"
	"github.com/jschalkwijk/GolangBlog/controller/categories"
	"path/filepath"
//	"runtime"
//	"path"
	"os"
	//"fmt"
	"fmt"
)

var static, _ = filepath.Abs("../jschalkwijk/GolangBlog/static/")

func main() {
// With this funtion I can check if my filepath is working for serving static files such as CSS or Templates etc
// IMPORTANT:I failed to add stat ic files because Go will use the current Directory you are in as the App's ROOT.
// If I run it from GolangBlog, the root is /Users/jorn/Documents/Golang/src/github.com/jschalkwijk/GolangBlog
// If I run it from jschalkwijk
	// 	_, err := os.Stat(filepath.Join(".", "GolangBlog/static/css", "style.css"))
//	checkErr(err)

	r := mux.NewRouter()
	// Index
	//#1
	r.PathPrefix("/GolangBlog/static/css").Handler(http.StripPrefix("/GolangBlog/static/css", http.FileServer(http.Dir("./GolangBlog/static/css"))))

	b,err := exists(static);
	fmt.Println(b)
	checkErr(err)
	r.HandleFunc("/", home.DashboardHandler)
	// Posts
	r.HandleFunc("/posts/", posts.Index)
		p := r.PathPrefix("/posts").Subrouter()
		p.HandleFunc("/{id:[0-9]+}/{title}", posts.Single)
		p.HandleFunc("/new", posts.New)
		p.HandleFunc("/edit/{id:[0-9]+}/{title}", posts.Edit)
		p.HandleFunc("/save/{id:[0-9]+}/{title}", posts.Save)
		p.HandleFunc("/add-post", posts.Add)
	// Categories
	r.HandleFunc("/categories/", categories.Index)
		c := r.PathPrefix("/categories").Subrouter()
		c.HandleFunc("/{id:[0-9]+}/{title}", categories.Single)
		c.HandleFunc("/new", categories.New)
		c.HandleFunc("/edit/{id:[0-9]+}/{title}", categories.Edit)
		c.HandleFunc("/save/{id:[0-9]+}/{title}", categories.Save)
		c.HandleFunc("/add-category", categories.Add)

//	s := http.FileServer(http.Dir(static))
//	r.PathPrefix("/static/").Handler(s)
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}


func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"css":           staticDirectory + "/css/",
		"images":           staticDirectory + "/images/",
		"scripts":          staticDirectory + "/scripts/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}