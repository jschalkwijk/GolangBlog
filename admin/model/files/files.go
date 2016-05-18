package files

import (
	"net/http"
	"html/template"
	cfg "github.com/jschalkwijk/GolangBlog/admin/config"
	"time"
	"fmt"
	"os"
	"io"
	"github.com/nu7hatch/gouuid"
	"strconv"
	"strings"
//	"database/sql"
//	_"database/sql/driver"
)

type File struct {
	FileID int
	Name string
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
	Messages []string
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

func insertRows(name string,fileName string,fileSize string,fileType string){
	fmt.Println("Name: ",name)
	fmt.Println("Filename: ",fileName)
	fmt.Println("Size: ", fileSize, " MB")
	fmt.Println("Type: ",fileType)

//	db, err := sql.Open("mysql", cfg.DB)
}
func Upload(w http.ResponseWriter, r *http.Request) {
	data := new(Data)
	// fmt.Println("method:", r.Method)

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20)
	checkErr(err)
	// Get a reference to the parsed multipart form
	m := r.MultipartForm
	// Access file headers
	files := m.File["uploadfile"]
//	file, handler, err := r.FormFile("uploadfile")
//	checkErr(err)
	for i, _ := range files {
		// For eah file header, get the handle to each file
		file, err := files[i].Open()
		defer file.Close()
		checkErr(err)

//		fmt.Fprintf(w, "%v", handler.Header)

		// Get the name, type and create a unique name to store in filesystem.
		name := files[i].Filename
		stripType := strings.Split(files[i].Header.Get("Content-Type"), "/")
		fileType := stripType[1]
		fileName := newName()

		// Check if files folder exists
		// if not create it.
		_, err = os.Stat("GolangBlog/files")
		if err != nil {
			err = os.Mkdir("GolangBlog/files", 0777)
			checkErr(err)
		}

		// Open a new empty file at a existing path plus the new file name and correct file typ
		f, err := os.OpenFile("GolangBlog/files/" + fileName + "." + fileType, os.O_WRONLY | os.O_CREATE, 0777)
		checkErr(err)
		defer f.Close()
		/*
			func Copy(dst Writer, src Reader) (written int64, err error)
			Copy copies from src to dst
		*/
		// Copy the uploaded file src to the new empty file on the filesystem.
		io.Copy(f, file)
		// Get the filesize of the file and convert to MB
		fileInfo, err := os.Stat("GolangBlog/files/" + fileName + "." + fileType)
		// bytes to MB. 1024 bytes = 1KB.
		fileSize := fmt.Sprintf("%0.2f", (float64(fileInfo.Size()) / 1024) / 1000)
		checkErr(err)

		// Insert into Database.
		insertRows(name, fileName, fileSize, fileType)
		data.Messages = append(data.Messages,name+" succesfully added to database.")
	}
	RenderTemplate(w,"files",data)

}

func newName() string {
	uuid, err := uuid.NewV4()
	checkErr(err)
	year,month,day := time.Now().Date()
	fmt.Println(year,month,day)
	fileName := month.String()+"-"+strconv.Itoa(day)+"-"+strconv.Itoa(year)+"-"+uuid.String()
	return fileName
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("OOPS something went wrong, you better fix it!", err)
		return
	}
}