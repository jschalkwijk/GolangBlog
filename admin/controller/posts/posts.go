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
	"github.com/jschalkwijk/GolangBlog/admin/model/categories"
	"github.com/gorilla/mux"
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"github.com/jschalkwijk/GolangBlog/admin/controller"
	"fmt"
)


func Index(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	fmt.Print(session.Username)
	fmt.Print(session.FirstName)
	fmt.Print(session.LastName)
	fmt.Print(session.Rights)
	fmt.Print(session.Logged)

	if (r.PostFormValue("approve-selected") != ""){
		a.Approve(w,r,"posts")
	}
	if (r.PostFormValue("trash-selected") != ""){
		a.Trash(w,r,"posts")
	}
	if (r.PostFormValue("hide-selected") != ""){
		a.Hide(w,r,"posts")
	}

	/** Example for interface implementation inside rendertemplate. **/
	//d := new(posts.Data)
	//p := d.GetPosts(0)

	p := posts.All(0)
	controller.RenderTemplate(w,"posts/posts", p)
}

func Deleted(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	if (r.PostFormValue("restore-selected") != ""){
		a.Restore(w,r,"posts")
	}
	if (r.PostFormValue("delete-selected") != ""){
		a.Delete(w,r,"posts")
	}
	p := posts.All(1)
	controller.RenderTemplate(w,"posts/posts", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	p := posts.Single(id,false)
	controller.RenderTemplate(w,"posts/posts", p)
}

func New(w http.ResponseWriter, r *http.Request){

	p,created := posts.Create(r);
	p.Categories = categories.GetCategories(0).Categories
	if(created){
		http.Redirect(w, r, "/admin/posts", http.StatusFound)
	} else {
		controller.RenderTemplate(w,"posts/add-edit-post",p)
	}
}

func Edit(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	p := posts.Single(id,true)

	if r.Method == "POST" {
		p, created := p.Posts[0].Patch(r);
		if (created) {
			http.Redirect(w, r, "/admin/posts", http.StatusFound)
		}  else {
			controller.RenderTemplate(w,"posts/add-edit-post",p)
		}
	} else {
		controller.RenderTemplate(w,"posts/add-edit-post",p)
	}
}