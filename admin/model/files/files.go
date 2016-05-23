package files

import (
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"net/http"
	"html/template"
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
	Type string
	Size string
	Path string
	Approved int
	Trashed int
	Date time.Time
	AlbumID int
}

type Data struct {
	Files []File
	Folders []Folder
	Deleted bool
	Messages []string
}

func RenderTemplate(w http.ResponseWriter, name string, f interface{}){
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

func Files(){

/*
	if(getCat) {
		listCat := cat.GetCategories(0)
		collection.Categories = listCat.Categories
	}
*/

}

func Upload(w http.ResponseWriter, r *http.Request) {
	data := new(Data)
	// fmt.Println("method:", r.Method)
	// Check if files folder exists
	// if not create it.
	_, err := os.Stat("GolangBlog/files")
	if err != nil {
		err = os.Mkdir("GolangBlog/files", 0777)
		checkErr(err)
	}

	folder := r.FormValue("new_folder_name")
	if (folder != ""){
		_, err = os.Stat("GolangBlog/files/"+folder)
		if err != nil {
			err = os.Mkdir("GolangBlog/files/"+folder, 0777)
			checkErr(err)
			folderPath := "files/"+folder
			folder := &Folder{FolderName: folder, ParentID: 0, Path: folderPath}
			err  = folder.create()
			checkErr(err)
		} else {
			fmt.Println("Folder already exists")
		}
	}
	// Parse multipart form
	err = r.ParseMultipartForm(32 << 20)
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
		path := "files/" + fileName + "." + fileType

		// Open a new empty file at a existing path plus the new file name and correct file typ
		f, err := os.OpenFile("GolangBlog/" + path, os.O_WRONLY | os.O_CREATE, 0777)
		checkErr(err)
		defer f.Close()
		/*
			func Copy(dst Writer, src Reader) (written int64, err error)
			Copy copies from src to dst
		*/
		// Copy the uploaded file src to the new empty file on the filesystem.
		io.Copy(f, file)
		// Get the filesize of the file and convert to MB
		fileInfo, err := os.Stat("GolangBlog/" + path)
		// bytes to MB. 1024 bytes = 1KB.
		fileSize := fmt.Sprintf("%0.2f", (float64(fileInfo.Size()) / 1024) / 1000)
		checkErr(err)

		// Assign values to the memory address of the file struct
		fl := &File{Name: name,FileName: fileName, Type: fileType,Size: fileSize,Path: path }
		// Call method insertRows to Insert into Database.
		err = fl.insertRows()
		data.Messages = append(data.Messages,name+" succesfully added to database.")
	}

	RenderTemplate(w,"files",data)

}

func (f *File) insertRows() error {
	fmt.Println("Name: ", f.Name)
	fmt.Println("Filename: ",f.FileName)
	fmt.Println("Size: ", f.Size, " MB")
	fmt.Println("Type: ",f.Type)

	db, err := sql.Open("mysql", config.DB)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO files (name,type,file_name,date,path) VALUES(?,?,?,FORMAT(Now(),'MM-DD-YYYY'),?)")
	fmt.Println(stmt)
	checkErr(err)

	res, err := stmt.Exec(f.Name,f.Type,f.FileName,f.Path)
	affect, err := res.RowsAffected()
	fmt.Println(affect)
	checkErr(err)

	return err
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