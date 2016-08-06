package files

import (
	"github.com/jschalkwijk/GolangBlog/admin/model/files"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"net/http"

)
func Index(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	f := files.Files()

	files.RenderTemplate(w,"files",f)
}

func Upload(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	println("hello");
	files.Upload(w,r)
}

