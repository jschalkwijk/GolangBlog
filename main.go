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

//	_, err := os.Stat(filepath.Join(".", static, "style.css"))
//	checkErr(err)

	r := mux.NewRouter()
	// Index
	//#1
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//#2
	//	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	//	r.PathPrefix("/static/").Handler(s)
	//#3
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//#4
	//ServeStatic(r, "/static/")
	//#5
	//r.Handle("/static/{rest}", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// print current dir
	//	_, filename, _, _ := runtime.Caller(1)
//	f, err := os.Open(path.Join(path.Dir(filename), ""))
//	fmt.Println(f)
//	checkErr(err)
	
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