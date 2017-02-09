package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	//front-end controllers
	"github.com/jschalkwijk/GolangBlog/controller/blog"
	cat "github.com/jschalkwijk/GolangBlog/controller/categories"
	"github.com/jschalkwijk/GolangBlog/model/home"
	//back-end controllers
	"github.com/jschalkwijk/GolangBlog/admin/controller/categories"
	"github.com/jschalkwijk/GolangBlog/admin/controller/dashboard"
	"github.com/jschalkwijk/GolangBlog/admin/controller/files"
	"github.com/jschalkwijk/GolangBlog/admin/controller/login"
	"github.com/jschalkwijk/GolangBlog/admin/controller/pages"
	"github.com/jschalkwijk/GolangBlog/admin/controller/posts"
	"github.com/jschalkwijk/GolangBlog/admin/controller/users"
)

var static string = "/static/"
var adminStatic string = "/admin/static/"

func main() {
	// With this function I can check if my filepath is working for serving static files such as CSS or Templates etc
	// IMPORTANT:I failed to add static files because Go will use the current Directory you are in as the App's ROOT.
	// If I run it from GolangBlog, the root is /Users/jorn/Documents/Golang/src/github.com/jschalkwijk/GolangBlog
	// If I run it from jschalkwijk
	// 	_, err := os.Stat(filepath.Join(".", "GolangBlog/static/css", "style.css"))
	//	checkErr(err)
	fmt.Println("Starting GolangBlog..")
	r := mux.NewRouter()
	serveStatic(r, static, "")
	serveStatic(r, adminStatic, "/admin")
	//r.PathPrefix("/admin/css/").Handler(http.StripPrefix("/admin/css/", http.FileServer(http.Dir("."+adminCSS))))

	r.HandleFunc("/", home.DashboardHandler)
		i := r.PathPrefix("/").Subrouter()
		i.HandleFunc("/{id:[0-9]+}/{title}", pages.One)
	// Blog
	r.HandleFunc("/blog", blog.Index)
		b := r.PathPrefix("/blog").Subrouter()
		b.HandleFunc("/{id:[0-9]+}/{title}", blog.One)
	// Categories
	r.HandleFunc("/categories", cat.Index)
		c := r.PathPrefix("/categories").Subrouter()
		c.HandleFunc("/{id:[0-9]+}/{title}", cat.One)
	//Admin
	r.HandleFunc("/admin", dashboard.Index)
	//Admin Posts
	r.HandleFunc("/admin/posts", posts.Index)
		aP := r.PathPrefix("/admin/posts/").Subrouter()
		aP.HandleFunc("/{id:[0-9]+}/{title}", posts.One)
		aP.HandleFunc("/add-post", posts.Create)
		aP.HandleFunc("/edit/{id:[0-9]+}/{title}", posts.Edit)
		aP.HandleFunc("/new", posts.Create)
		aP.HandleFunc("/trashed-posts", posts.Deleted)
	// Admin Pages
	r.HandleFunc("/admin/pages", pages.Index)
		aPa := r.PathPrefix("/admin/pages/").Subrouter()
		//pages.HandleFunc("/trashed-pages", pages.Deleted)
		aPa.HandleFunc("/{id:[0-9]+}/{title}", pages.One)
		aPa.HandleFunc("/new", pages.Create)
		aPa.HandleFunc("/edit/{id:[0-9]+}", pages.Edit)

	//Admin Categories
	r.HandleFunc("/admin/categories", categories.Index)
		aC := r.PathPrefix("/admin/categories").Subrouter()
		aC.HandleFunc("/{id:[0-9]+}/{title}", categories.One)
		aC.HandleFunc("/add-category", categories.Create)
		aC.HandleFunc("/edit/{id:[0-9]+}/{title}", categories.Edit)
		aC.HandleFunc("/add", categories.Add)
		aC.HandleFunc("/trashed-categories", categories.Deleted)
	// Users
	r.HandleFunc("/admin/users", users.Index)
		u := r.PathPrefix("/admin/users").Subrouter()
		u.HandleFunc("/{id:[0-9]+}/{username}", users.One)
		u.HandleFunc("/add-user", users.Create)
		u.HandleFunc("/add", users.Add)
		u.HandleFunc("/edit/{id:[0-9]+}/{username}", users.Edit)
		u.HandleFunc("/trashed-users", users.Deleted)
	//Files
	r.HandleFunc("/admin/files", files.Index)
	f := r.PathPrefix("/admin/files/").Subrouter()
	f.HandleFunc("/upload", files.Upload)
	f.HandleFunc("/folder/{id:[0-9]+}/{foldername}", files.Folder)
	// Login
	r.HandleFunc("/admin/login", login.Index)
		l := r.PathPrefix("/admin/login").Subrouter()
		l.HandleFunc("/auth", login.Auth)
		l.HandleFunc("/logout", login.Logout)

	http.Handle("/", r)

	fmt.Println("GolangBlog running on port 8080. Don't forget to run MAMP or SQL server.")

	http.ListenAndServe(":8080", nil)
}

// load all static directories: Source: http://www.shakedos.com/2014/Feb/08/serving-static-files-with-go.html
func serveStatic(router *mux.Router, staticDirectory string, admin string) {
	staticPaths := map[string]string{
		"/css/":     staticDirectory + "/css/",
		"/test/":    staticDirectory + "/test/",
		"/images/":  staticDirectory + "/images/",
		"/scripts/": staticDirectory + "/scripts/",
		"/tinymce/": staticDirectory + "/scripts/tinymce/js/tinymce/",
		// If we use "/files/" as a prefix we get in conflict with the router which also use files.
		// Also it only works if the files folder is inside another folder also due to the conflict.
		"/uploads/": staticDirectory + "/uploads/",
	}
	for pathPrefix, pathValue := range staticPaths {
		router.PathPrefix(admin + pathPrefix).Handler(http.StripPrefix(admin+pathPrefix,
			http.FileServer(http.Dir("."+pathValue))))
	}
}
