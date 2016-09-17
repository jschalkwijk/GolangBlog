package files

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/jschalkwijk/GolangBlog/admin/model/files"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
)
func Index(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
	f := files.Files("","")

	if r.PostFormValue("action") == "trash" {
		// Form submitted
		a.Trash(w,r,"files")
	}
	if r.PostFormValue("action") == "restore" {
		// Form submitted
		a.Restore(w,r,"files")
	}
	if r.PostFormValue("action") == "delete" {
		// Form submitted
		a.Delete(w,r,"files")
	}
	if r.PostFormValue("action") == "delete-folder" {
		// Form submitted
		msg := a.DeleteFolders(w,r,"files")
		f.Messages = msg
	}

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
	idINT, _ := strconv.Atoi(id)

	f.CurrentFolder = idINT

	files.RenderTemplate(w,"files",f)
}
