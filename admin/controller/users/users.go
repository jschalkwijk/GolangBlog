package users

import (
	"net/http"
	 a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/users"

)

var dbt string = "users"

func Index(w http.ResponseWriter, r *http.Request) {
	if (r.PostFormValue("approve-selected") != ""){
		a.Approve(w,r,dbt)
	}
	if (r.PostFormValue("trash-selected") != ""){
		a.Trash(w,r,dbt)
	}
	if (r.PostFormValue("hide-selected") != ""){
		a.Hide(w,r,dbt)
	}
	data := new(users.Data)
	u := data
	users.RenderTemplate(w,"users", u)
}
