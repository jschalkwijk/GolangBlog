package posts

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/model/posts"
	"github.com/jschalkwijk/GolangBlog/controller"
	"github.com/gorilla/mux"
)


func Index(w http.ResponseWriter, r *http.Request) {
	p := posts.GetPosts()
	controller.RenderTemplate(w,"posts", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := posts.GetSinglePost(id,post_title)
	controller.RenderTemplate(w,"posts", p)
}

func New(w http.ResponseWriter, r *http.Request){
	collection := new(posts.Data)
	p := collection
	controller.RenderTemplate(w,"new-post", p)
}

func Edit(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := posts.GetSinglePost(id,post_title)
	controller.RenderTemplate(w,"edit-post", p)
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


