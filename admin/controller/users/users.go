package users

import (
	"net/http"
	 a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/users"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"github.com/gorilla/mux"
	"github.com/jschalkwijk/GolangBlog/admin/controller"

)

var dbt string = "users"

func Index(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	if (r.PostFormValue("approve-selected") != ""){
		a.Approve(w,r,dbt)
	}
	if (r.PostFormValue("trash-selected") != ""){
		a.Trash(w,r,dbt)
	}
	if (r.PostFormValue("hide-selected") != ""){
		a.Hide(w,r,dbt)
	}
	u := users.All(0)
	controller.RenderTemplate(w,"users/users", u)
}

func Deleted(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	if (r.PostFormValue("restore-selected") != ""){
		a.Restore(w,r,dbt)
	}
	if (r.PostFormValue("delete-selected") != ""){
		a.Delete(w,r,dbt)
	}
	p := users.All(1)
	controller.RenderTemplate(w,"users/users", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	p := users.Single(id)
	controller.RenderTemplate(w,"users/users", p)
}

func New(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	u,created := users.Create(r)

	if(created){
		http.Redirect(w, r, "/admin/users", http.StatusFound)
	} else {
		controller.RenderTemplate(w,"users/add-user",u)
	}
}

func Edit(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	u := users.Single(id)

	if r.Method == "POST" {
		_, updated := u.Users[0].Patch(r);
		if (updated) {
			http.Redirect(w, r, "/admin/users", http.StatusFound)
		}
	}
	controller.RenderTemplate(w,"users/edit-user",u)
}

func Add(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	u,created := users.Create(r)

	if(created){
		http.Redirect(w, r, "/admin/users", http.StatusFound)
	} else {
		controller.RenderTemplate(w,"users/add-user",u)
	}
}
