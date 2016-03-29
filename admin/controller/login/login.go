package login

import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	//"github.com/gorilla/mux"
)



func Index(w http.ResponseWriter, r *http.Request) {
	login.RenderTemplate(w,"login")
}