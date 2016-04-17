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
	data := new(files.Data)

	files.RenderTemplate(w,"files",data)
}
