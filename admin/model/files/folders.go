package files

import (
	"fmt"
	"os"
	"database/sql"
	_"database/sql/driver"
	"github.com/jschalkwijk/GolangBlog/admin/config"
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
}

var folderID int
var folderName string
var description string
var author string
var parentID int
var folderPath string
var date string

func Folders() []Folder {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	rows, err := db.Query("SELECT * FROM folders")
	checkErr(err)

	var folders []Folder
	for rows.Next() {
		err = rows.Scan(&folderID, &folderName,&description,&author,&parentID,&folderPath,&date)
		checkErr(err)

		// calculate folder size. path.WALK functie proberen, hij moet de grote van alle bestanden bij elkaar optellen
		// dat gebeurd nu niet. check http://stackoverflow.com/questions/32482673/golang-how-to-get-directory-total-size
		info, err := os.Lstat("GolangBlog/static/"+folderPath)
		checkErr(err)
		println("Folder-Size",info.Size())
		sizeMB := fmt.Sprintf("%0.2f", (float64(info.Size()) / 1024) / 1024)
		println(sizeMB)
		//

		folder := Folder{folderID, folderName,description,author,parentID,folderPath,date}
		folders = append(folders, folder)
	}
	return folders
}

func (folder *Folder) save() error{
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	stmt, err := db.Prepare("INSERT INTO folders(folder_name,path) VALUES(?,?)")
	fmt.Println(stmt)
	checkErr(err)
	_, err = stmt.Exec(folder.FolderName, folder.FolderPath)
	checkErr(err)
	return err
}

func Create(folder string) error {
	// Check if files folder exists
	// if not create it.
	_, err := os.Stat("GolangBlog/static/files/"+folder)
	if err != nil {
		err = os.Mkdir("GolangBlog/static/files/"+folder, 0777)
		checkErr(err)
	}

	newFolder := Folder{FolderName: folder,FolderPath: "files/"+folder}
	err = newFolder.save();
	checkErr(err)
	return err;
}

