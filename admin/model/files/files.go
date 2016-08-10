package files

import (
	"net/http"
	"html/template"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"time"
	"fmt"
	"os"
	"io"
	"github.com/nu7hatch/gouuid"
	"strconv"
	"strings"
	"database/sql"
	_"database/sql/driver"
)

type File struct {
	FileID int
	Name string
	FileName string
	FileType string
	Size string
	FilePath string
	//Location string
	//Approved int
	//Trashed int
	//Date time.Time
	FolderID int
}

type Data struct {
	Files []File
	Folders []Folder
	CurrentFolder int
	Deleted bool
	Messages []string
}

var file_id int
var name string
var fileName string
var fileType string
var size string
var filePath string
var fID int
func RenderTemplate(w http.ResponseWriter, name string, f *Data){
	t, err := template.ParseFiles(config.Templates+"/"+"header.html",config.Templates+"/"+"nav.html",config.View + "/" + name + ".html",config.Templates+"/"+"footer.html")
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

func Files(folderID string,folderName string) *Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	rows, err := db.Query("SELECT file_id,name,file_name,type,size,path,folder_id FROM files WHERE folder_id = ? ORDER BY file_id DESC",folderID)
	checkErr(err)

	data := new(Data)

	for rows.Next() {
		err = rows.Scan(&file_id, &name, &fileName,&fileType, &size, &filePath,&fID)
		checkErr(err)
		//fID,err = strconv.Atoi(fID)
		//checkErr(err)
		file := File{file_id, name, fileName, fileType, size, filePath, fID}
		data.Files = append(data.Files, file)
	}
	 data.Folders = Folders(folderID)

	return data
}

func (f *File) insertRows() error {
	fmt.Println("Name: ",f.Name)
	fmt.Println("Filename: ",f.FileName)
	fmt.Println("Size: ", f.Size, " MB")
	fmt.Println("Type: ",f.FileType)
	fmt.Println("Path: ",f.FilePath)
	fmt.Println("FolderID: ",f.FolderID)
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO files (name,file_name,size,type,path,folder_id) VALUES(?,?,?,?,?,?)")
	fmt.Println(stmt)
	checkErr(err)
	_, err = stmt.Exec(f.Name, f.FileName,f.Size,f.FileType,f.FilePath,f.FolderID)
	checkErr(err)
	return err
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
	folder := r.FormValue("new_folder_name")
	parentID := r.FormValue("folder_name")
	var lastID int = 0
	if (folder != "") {
		lastID, err = Create(folder,parentID)
		checkErr(err)
	}

	fmt.Println("LastID: ", lastID)

	for i, _ := range files {
		// For eah file header, get the handle to each file
		file, err := files[i].Open()
		defer file.Close()
		checkErr(err)

		//		fmt.Fprintf(w, "%v", handler.Header)

		// Get the name, type and create a unique name to store in filesystem.
		fname := files[i].Filename
		stripType := strings.Split(files[i].Header.Get("Content-Type"), "/")
		fType := stripType[1]
		fName := newName()
		// If we use "/files/" as a prefix we get in conflict with the router which also use files.
		// Also it only works if the files folder is inside another folder also due to the conflict.
		// see main.go.
		fPath:= "/file/" + fName + "." + fType

		// Check if files folder exists
		// if not create it.
		_, err = os.Stat("GolangBlog/static/files")
		if err != nil {
			err = os.Mkdir("GolangBlog/static/files", 0777)
			checkErr(err)
		}

		// Open a new empty file at a existing path plus the new file name and correct file typ
		f, err := os.OpenFile("GolangBlog/static/files/" + fName + "." + fType, os.O_WRONLY | os.O_CREATE, 0777)
		checkErr(err)
		defer f.Close()
		/*
			func Copy(dst Writer, src Reader) (written int64, err error)
			Copy copies from src to dst
		*/
		// Copy the uploaded file src to the new empty file on the filesystem.
		io.Copy(f, file)
		// Get the filesize of the file and convert to MB.
		// Stat returns a FileInfo structure describing the named file.
		fileInfo, err := os.Stat("GolangBlog/static/files/" + fName + "." + fType)
		// bytes to MB. 1024 bytes = 1KB.
		fSize := fmt.Sprintf("%0.2f", (float64(fileInfo.Size()) / 1024) / 1024)
		checkErr(err)

		media := &File{ Name: fname,FileName: fName,Size: fSize,FileType: fType,FilePath: fPath,FolderID:lastID }
		// Insert into Database.
		err = media.insertRows()
		checkErr(err)

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
		fmt.Println("OOPS something went wrong in the files model, you better fix it!", err)
		return
	}
}