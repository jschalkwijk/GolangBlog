package posts

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/posts"
	//"github.com/jschalkwijk/GolangBlog/controller"
	"github.com/gorilla/mux"
	//"github.com/jschalkwijk/GolangBlog/model/data"
	"github.com/jschalkwijk/GolangBlog/admin/model/categories"
)


func Index(w http.ResponseWriter, r *http.Request) {
	p := posts.GetPosts()
	posts.RenderTemplate(w,"posts", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := posts.GetSinglePost(id,post_title,false)
	posts.RenderTemplate(w,"posts", p)
}

func New(w http.ResponseWriter, r *http.Request){
	c := categories.GetCategories()
	categories.RenderTemplate(w,"add-post", c)
}

func Edit(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := posts.GetSinglePost(id,post_title, true)
	posts.RenderTemplate(w,"edit-post", p)
}
func Save(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	posts.EditPost(w,r,id,post_title)
}

func Add(w http.ResponseWriter, r *http.Request){
	posts.NewPost(w, r)
}


