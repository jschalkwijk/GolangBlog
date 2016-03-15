/*	-- Categories Controller --
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

package categories

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/categories"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	p := categories.GetCategories()
	categories.RenderTemplate(w,"categories", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := categories.GetSingleCategory(id,post_title)
	categories.RenderTemplate(w,"categories", p)
}

func New(w http.ResponseWriter, r *http.Request){
	collection := new(categories.Data)
	p := collection
	categories.RenderTemplate(w,"add-category", p)
}

func Edit(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := categories.GetSingleCategory(id,post_title)
	categories.RenderTemplate(w,"edit-category", p)
}

func Save(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	categories.EditCategory(w,r,id)
}

func Add(w http.ResponseWriter, r *http.Request){
	categories.NewCategory(w, r)
}