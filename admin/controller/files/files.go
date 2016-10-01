package files

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/jschalkwijk/GolangBlog/admin/model/files"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/controller"
)
func Index(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	files := new(files.Data)
	f := files.Get("","")

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

	controller.RenderTemplate(w,"files",f)
}

func Upload(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	files.Upload(w,r)
	http.Redirect(w,r,"/admin/files",http.StatusFound);
}

func Folder(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w,r, "/admin/login", http.StatusFound)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	folderName := vars["foldername"]

	files := new(files.Data)
	f := files.Get(id,folderName)

	idINT, _ := strconv.Atoi(id)

	f.CurrentFolder = idINT

	controller.RenderTemplate(w,"files",f)
}
