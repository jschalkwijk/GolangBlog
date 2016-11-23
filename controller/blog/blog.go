package blog

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/model/blog"
	"github.com/jschalkwijk/GolangBlog/admin/model/posts"
	//"github.com/jschalkwijk/GolangBlog/controller"
	"github.com/gorilla/mux"
	//"github.com/jschalkwijk/GolangBlog/model/data"
)


func Index(w http.ResponseWriter, r *http.Request) {
	p := posts.All(0)
	blog.RenderTemplate(w,"blog", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	p := posts.Single(id,false)
	blog.RenderTemplate(w,"blog", p)
}

