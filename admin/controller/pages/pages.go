package pages

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"github.com/jschalkwijk/GolangBlog/admin/model/pages"
	"github.com/jschalkwijk/GolangBlog/admin/controller"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	p := pages.All(0)
	controller.View(w,"pages/pages",p)
}

func One(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	p := pages.One(id)
	controller.View(w,"pages/pages",p)
}

func Create(w http.ResponseWriter, r *http.Request){

	p,created := pages.Create(r);
	if(created){
		http.Redirect(w, r, "/admin/pages", http.StatusFound)
	}

	controller.View(w,"pages/add-edit-page",p)
}

func Edit(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	p := pages.One(id)

	if r.Method == "POST" {
		data, created := p.Pages[0].Patch(r);
		if (created) {
			http.Redirect(w, r, "/admin/pages", http.StatusFound)
		}  else {
			p = data
		}
	}
	controller.View(w,"pages/add-edit-page",p)
}