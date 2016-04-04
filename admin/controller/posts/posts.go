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
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"fmt"
)


func Index(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (session.Logged != true) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	fmt.Printf("%s", session.Username)
	fmt.Printf("%s", session.FirstName)
	fmt.Printf("%s", session.LastName)
	fmt.Printf("%s", session.Rights)
	fmt.Printf("%s", session.Logged)

	if (r.PostFormValue("approve-selected") != ""){
		a.Approve(w,r,"posts")
	}
	if (r.PostFormValue("trash-selected") != ""){
		a.Trash(w,r,"posts")
	}
	if (r.PostFormValue("hide-selected") != ""){
		a.Hide(w,r,"posts")
	}
	p := posts.GetPosts(0)
	posts.RenderTemplate(w,"posts", p)
}

func Deleted(w http.ResponseWriter, r *http.Request) {
	if (r.PostFormValue("restore-selected") != ""){
		a.Restore(w,r,"posts")
	}
	if (r.PostFormValue("delete-selected") != ""){
		a.Delete(w,r,"posts")
	}
	p := posts.GetPosts(1)
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
	c := categories.GetCategories(0)
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


