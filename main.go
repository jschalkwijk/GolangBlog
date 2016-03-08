package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	//front-end controllers
	"github.com/jschalkwijk/GolangBlog/model/home"
	"github.com/jschalkwijk/GolangBlog/controller/blog"
	cat "github.com/jschalkwijk/GolangBlog/controller/categories"
	//back-end controllers
	"github.com/jschalkwijk/GolangBlog/admin/model/admin"
	"github.com/jschalkwijk/GolangBlog/admin/controller/posts"
	"github.com/jschalkwijk/GolangBlog/admin/controller/categories"
)

var static string = "/GolangBlog/static/"
var adminStatic string = "/GolangBlog/admin/static/"

func main() {
	// With this funtion I can check if my filepath is working for serving static files such as CSS or Templates etc
	// IMPORTANT:I failed to add static files because Go will use the current Directory you are in as the App's ROOT.
	// If I run it from GolangBlog, the root is /Users/jorn/Documents/Golang/src/github.com/jschalkwijk/GolangBlog
	// If I run it from jschalkwijk
		// 	_, err := os.Stat(filepath.Join(".", "GolangBlog/static/css", "style.css"))
	//	checkErr(err)
	fmt.Println("Starting GolangBlog..")
	r := mux.NewRouter()
	serveStatic(r,static,"")
	serveStatic(r,adminStatic,"/admin")
	//r.PathPrefix("/admin/css/").Handler(http.StripPrefix("/admin/css/", http.FileServer(http.Dir("."+adminCSS))))

	r.HandleFunc("/", home.DashboardHandler)
	// Blog
	r.HandleFunc("/blog", blog.Index)
		b := r.PathPrefix("/blog").Subrouter()
		b.HandleFunc("/{id:[0-9]+}/{title}", blog.Single)
	// Categories
	r.HandleFunc("/categories", cat.Index)
		c := r.PathPrefix("/categories").Subrouter()
		c.HandleFunc("/{id:[0-9]+}/{title}", cat.Single)
	//Admin
	r.HandleFunc("/admin", admin.DashboardHandler)
	//Admin Posts
	r.HandleFunc("/admin/posts", posts.Index)
		aP := r.PathPrefix("/admin/posts/").Subrouter()
		aP.HandleFunc("/{id:[0-9]+}/{title}", posts.Single)
		aP.HandleFunc("/add-post", posts.New)
		aP.HandleFunc("/edit/{id:[0-9]+}/{title}", posts.Edit)
		aP.HandleFunc("/save/{id:[0-9]+}/{title}", posts.Save)
		aP.HandleFunc("/add", posts.Add)
	//Admin Categories
	r.HandleFunc("/admin/categories", categories.Index)
		aC := r.PathPrefix("/admin/categories").Subrouter()
		aC.HandleFunc("/{id:[0-9]+}/{title}", categories.Single)
		aC.HandleFunc("/add-category", categories.New)
		aC.HandleFunc("/edit/{id:[0-9]+}/{title}", categories.Edit)
		aC.HandleFunc("/save/{id:[0-9]+}/{title}", categories.Save)
		aC.HandleFunc("/add", categories.Add)

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
func serveStatic(router *mux.Router, staticDirectory string, admin string) {
	staticPaths := map[string]string{
		"/css/":           staticDirectory + "/css/",
		"/images/":           staticDirectory + "/images/",
		"/scripts/":          staticDirectory + "/scripts/",
	}
	for pathPrefix, pathValue := range staticPaths {
		router.PathPrefix(admin+pathPrefix).Handler(http.StripPrefix(admin+pathPrefix,
			http.FileServer(http.Dir("."+pathValue))))
	}
}