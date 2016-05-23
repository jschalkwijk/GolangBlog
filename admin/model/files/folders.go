package files

import (
	"database/sql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
)

type Folder struct {
	FolderID int
	FolderName string
	Description string
	Author string
	ParentID int
	Path string
}

type FoldersData struct {
	Folders []Folder
}

var folderID int
var folderName string
var description string
var author string
var parentID int
var path string



func (f *Folder) create() error {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	query , err := db.Prepare("INSERT INTO folders (folder_name,parent_id,path) VALUES(?,?,?)")
	checkErr(err)

	_ , err = query.Exec(f.FolderName,f.ParentID,f.Path)
	checkErr(err)

	return err
}

func Folders(parentID int) *FoldersData {
	db,err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM folders WHERE parent_id = ?", parentID)
	checkErr(err)

	data := new(FoldersData)

	for rows.Next() {
		err = rows.Scan(&folderID, &folderName, &description,&author,&parentID, &path)
		checkErr(err)
		// convert string to HTML markdown
		folder := Folder{folderID,folderName,description,author,parentID,path}
		data.Folders = append(data.Folders , folder)
	}

	return data
}