package files

import (

)

type Folder struct {
	FolderID int
	FolderName string
	Path string
	Approved int
	Trashed int
	ParentID int
}

func (f *Folder) create() {

}