package login

import (
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	//"github.com/gorilla/mux"

)

func Index(w http.ResponseWriter, r *http.Request) {
	login.View(w,"login")
}

func Auth(w http.ResponseWriter, r *http.Request) {
	auth := login.Login(w, r)
	if (auth == nil) {
		login.SetSession(w,r)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		login.View(w,"login")
	}
}

func Logout(w http.ResponseWriter, r *http.Request){
	login.Logout(w,r)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

