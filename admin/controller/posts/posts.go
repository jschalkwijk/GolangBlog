/*	-- PostsController --
 * 	All functions in this file are called by package main.
 *	Functions inside this controller take a http.ResponseWriter, r *http.Request.
 *  If specified in main we can take URL parameters using the Gorrila Mux tool.
 * 	They can call functions from an imported model.
 * 	If a func from a model returns data, it had to be assigned to a variable.
 *  The variable with the data must be passed to the models RenderTemplate func
 	in order to render the template with the data.
 *	In some cases you need to render a template without data. This is done by
 	creating an empty data struct from the imported model and then pass it to the
 	RenderTemplate func.
 */

package posts

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/posts"
	"github.com/gorilla/mux"
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


