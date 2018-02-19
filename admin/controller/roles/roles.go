package roles

import (
	"net/http"
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"github.com/gorilla/mux"
	"github.com/jschalkwijk/GolangBlog/admin/controller"

	"github.com/jschalkwijk/GolangBlog/admin/model/roles"
)

var dbt string = "roles"

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
	role := roles.All()
	controller.View(w,"roles/roles", role)
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
	p := roles.All()
	controller.View(w,"roles/roles", p)
}

func One(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	p := roles.One(id,false)
	controller.View(w,"roles/role", p)
}

//func Create(w http.ResponseWriter, r *http.Request){
//	session := login.GetSession(r)
//
//	if (!session.Logged) {
//		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
//	}
//
//	u,created := roles.Create(r)
//
//	if(created){
//		http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
//	} else {
//		controller.View(w,"roles/add-user",u)
//	}
//}
//
//func Edit(w http.ResponseWriter, r *http.Request){
//	session := login.GetSession(r)
//
//	if (!session.Logged) {
//		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
//	}
//
//	vars := mux.Vars(r)
//	id := vars["id"]
//	u := roles.One(id)
//
//	if r.Method == "POST" {
//		data, updated := u.Roles[0].Patch(r);
//		if (updated) {
//			http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
//		} else {
//			u = data
//		}
//	}
//	controller.View(w,"roles/edit-user",u)
//}
//
//func Add(w http.ResponseWriter, r *http.Request){
//	session := login.GetSession(r)
//
//	if (!session.Logged) {
//		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
//	}
//
//	u,created := roles.Create(r)
//
//	if(created){
//		http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
//	} else {
//		controller.View(w,"roles/add-user",u)
//	}
//}

