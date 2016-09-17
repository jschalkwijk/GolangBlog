package files

import (
	"fmt"
	"os"
	"database/sql"
	_"database/sql/driver"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"path/filepath"
)
type Folder struct {
	FolderID int
	FolderName string
	Description string
	//	Approved int
	//	Trashed int
	Author string
	ParentID int
	FolderPath string
	Date string
	FolderSize string
}

var folderID int
var folderName string
var description string
var author string
var parentID int
var folderPath string
var date string

func Folders(id string) []Folder {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	var rows *sql.Rows
	if(id == "") {
		rows, err = db.Query("SELECT * FROM folders WHERE parent_id = ? ORDER BY folder_id DESC",0)
	} else {
		rows, err = db.Query("SELECT * FROM folders WHERE folder_id = ? OR parent_id = ? ORDER BY folder_id DESC",id,id)
	}
	checkErr(err)

	var folders []Folder
	for rows.Next() {
		err = rows.Scan(&folderID, &folderName,&description,&author,&parentID,&folderPath,&date)
		checkErr(err)

		// calculate folder size. path.WALK functie proberen, hij moet de grote van alle bestanden bij elkaar optellen
		// dat gebeurd nu niet. check http://stackoverflow.com/questions/32482673/golang-how-to-get-directory-total-size
		size,err := DirSize("GolangBlog/static/"+folderPath)
		checkErr(err)

		sizeMB := fmt.Sprintf("%0.2f",size)
		fmt.Println(sizeMB)
		folder := Folder{folderID, folderName,description,author,parentID,folderPath,date,sizeMB}
		folders = append(folders, folder)
	}

	return folders
}
func (folder *Folder) save() (int, error){
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	stmt, err := db.Prepare("INSERT INTO folders(folder_name,path,parent_id) VALUES(?,?,?)")
	fmt.Println(stmt)
	checkErr(err)
	result, err := stmt.Exec(folder.FolderName, folder.FolderPath,folder.ParentID)
	checkErr(err)
	lastID, err := result.LastInsertId()
	checkErr(err)
	return int(lastID),err
}

func Create(folder string,parentID int) (int,string,error) {
	// Check if files folder exists
	// if not create it.
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)

	checkErr(err)
	var path string
	var row *sql.Row
	if(parentID != 0) {
		row = db.QueryRow("SELECT path FROM folders WHERE folder_id = ?", parentID)
		row.Scan(&path)
		path = path+"/"+folder
	} else {
		// with no parent folder we need to add the files/ prefix
		path = "files/"+folder
	}

	_, err = os.Stat("GolangBlog/static/"+path)
	if err != nil {
		err = os.Mkdir("GolangBlog/static/"+path, 0777)
		checkErr(err)
	}

	newFolder := Folder{FolderName: folder,FolderPath: path,ParentID: parentID}
	lastID, err := newFolder.save();
	checkErr(err)
	return lastID, path, err;
}

//func Remove (path string)(msg string){
//	var msg string
//	if(os.Stat(path)) {
//		err := os.RemoveAll(path)
//		checkErr(err)
//		msg = path + "is removed successfully"
//		return msg
//	} else {
//		msg = "The folder you want to delete doesn't exist"
//		return msg
//	}
//}

//this should be done only when the folder changes and then store into the DB
// not everytime we load the page.
func DirSize(path string)(float64, error){
	var size float64
	//Walk walks the file tree from the given filepath or root
	// Using a closure we can get the fileinfo and size of each file which will be appended to the size var.
	// returns the size in bites and a error message.
	err := filepath.Walk("GolangBlog/static/"+folderPath, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += (float64(info.Size()) / 1024) / 1024

			fmt.Println("foldersize: ",size)
		}
		return err
	})
	return size,err
}

