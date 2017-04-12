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
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
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
	controller.View(w,"users/users", u)
}

func Deleted(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	if (r.PostFormValue("restore-selected") != ""){
		a.Restore(w,r,dbt)
	}
	if (r.PostFormValue("delete-selected") != ""){
		a.Delete(w,r,dbt)
	}
	p := users.All(1)
	controller.View(w,"users/users", p)
}

func One(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	p := users.One(id)
	controller.View(w,"users/users", p)
}

func Create(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	u,created := users.Create(r)

	if(created){
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	} else {
		controller.View(w,"users/add-user",u)
	}
}

func Edit(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	u := users.One(id)

	if r.Method == "POST" {
		data, updated := u.Users[0].Patch(r);
		if (updated) {
			http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		} else {
			u = data
		}
	}
	controller.View(w,"users/edit-user",u)
}

func Add(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	u,created := users.Create(r)

	if(created){
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	} else {
		controller.View(w,"users/add-user",u)
	}
}
