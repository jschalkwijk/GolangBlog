package files

import (
	"github.com/jschalkwijk/GolangBlog/admin/model/files"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
)
func Index(w http.ResponseWriter, r *http.Request){
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	data := new(files.Data)
	folders := files.Folders(0)
	data.Folders = folders.Folders

	files.RenderTemplate(w,"files",data)
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
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	parentID, err := strconv.Atoi(id)
	checkErr(err)

	data := new(files.Data)
	folders := files.Folders(parentID)
	data.Folders = folders.Folders

	files.RenderTemplate(w,"files",data)

}

func checkErr(err error) {
	if err != nil {
		fmt.Println("OOPS something went wrong, you better fix it!", err)
		return
	}
}