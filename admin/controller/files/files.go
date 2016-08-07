package files

import (
	"github.com/jschalkwijk/GolangBlog/admin/model/files"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"net/http"
	"github.com/gorilla/mux"

)
func Index(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	f := files.Files("","")

	files.RenderTemplate(w,"files",f)
}

func Upload(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	files.Upload(w,r)
}

func Folder(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w,r, "/admin/login", http.StatusFound)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	folderName := vars["foldername"]
	f := files.Files(id,folderName)

	files.RenderTemplate(w,"files",f)
}

