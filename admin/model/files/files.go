package files

import (
	"net/http"
	"html/template"
	cfg "github.com/jschalkwijk/GolangBlog/admin/config"
	"time"
	"fmt"
	"os"
	"strconv"
	"crypto/md5"
	"io"
)

type File struct {
	FileID int
	FileName string
	Ext string
	Location string
	Approved int
	Trashed int
	Date time.Time
	AlbumID int
}

type Album struct {
	AlbumID int
	AlbumName string
	Ext string
	Location string
	Approved int
	Trashed int
	Date time.Time
	ParentID int
}

type Data struct {
	Files []File
	Albums []Album
	Deleted bool
}

func RenderTemplate(w http.ResponseWriter, name string, f *Data){
	t, err := template.ParseFiles(cfg.Templates+"/"+"header.html",cfg.Templates+"/"+"nav.html",cfg.View + "/" + name + ".html",cfg.Templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,f)
	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Files(){

}
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}