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
	d := new(users.Data)
	u := d.Get(0)
	controller.RenderTemplate(w,"users", u)
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
	d := new(users.Data)
	u := d.Get(0)
	controller.RenderTemplate(w,"users", u)
}

func Single(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	username := vars["username"]
	p := users.GetSingleUser(id,username)
	controller.RenderTemplate(w,"users", p)
}

func New(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	data := new(users.Data)
	u := data
	controller.RenderTemplate(w,"add-user", u)
}

func Edit(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	username := vars["username"]
	u := users.GetSingleUser(id,username)
	controller.RenderTemplate(w,"edit-user", u)
}

func Save(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	username := vars["username"]
	users.EditUser(w,r,id,username)
}

func Add(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	users.NewUser(w, r)
}
