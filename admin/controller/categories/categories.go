/*	-- Categories Controller --
 * 	All functions in this file are called by package main.
 *	Functions inside this controller take a http.ResponseWriter, r *http.Request.
 *  If specified in main we can take URL parameters using the Gorrila Mux tool.
 * 	They can call functions from an imported model.
 * 	If a func from a model returns data, it had to be assigned to a variable.
 *  The variable with the data must be passed to the models View func
 	in order to render the template with the data.
 *	In some cases you need to render a template without data. This is done by
 	creating an empty data struct from the imported model and then pass it to the
 	View func.
 */

package categories

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/controller"
	"github.com/jschalkwijk/GolangBlog/admin/model/categories"
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"github.com/gorilla/mux"
)

var dbt string = "categories"

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
	p := categories.All(0)
	controller.View(w,"categories/categories", p)
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
	p := categories.All(1)
	controller.View(w,"categories/categories", p)
}
func One(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	p := categories.One(id)
	controller.View(w,"categories/categories", p)
}

func Create(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	collection := new(categories.Data)
	p := collection
	controller.View(w,"categories/add-category", p)
}

func Edit(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	c := categories.One(id)

	if r.Method == "POST" {
		data, updated := c.Categories[0].Patch(r);
		if (updated) {
			http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
		} else {
			c = data
		}
	}
	controller.View(w,"categories/edit-category",c)
}

func Add(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	u,created := categories.Create(r)

	if(created){
		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
	} else {
		controller.View(w,"categories/add-category",u)
	}
}