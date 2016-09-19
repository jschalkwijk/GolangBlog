package actions

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"fmt"
	"strings"
	"os"
)

func Trash(w http.ResponseWriter, r *http.Request, dbt string){
	update(w,r, dbt,"trashed",1)
}
func Restore(w http.ResponseWriter, r *http.Request, dbt string){
	update(w,r, dbt,"trashed",0)
}
func Approve(w http.ResponseWriter, r *http.Request, dbt string){
	update(w,r, dbt,"approved",1)
}
func Hide(w http.ResponseWriter, r *http.Request, dbt string){
	update(w,r, dbt,"approved",0)
}

func update(w http.ResponseWriter, r *http.Request, dbt string,row string,setTo int) {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)
	// Getting the selected checkboxes from the request.Form
	checked := r.Form["checkbox"]
	// !!! fix a way to use the multiple function and add the setTo value to the map
	// something like args = append([]string{strconv.Itoa(setTo)},args)

	// Creating a string with the amount of ? required for the Query string by using the length of the checkbox slice.
	multiple := strings.Repeat("?, ",len(checked))
	// delete the last 2 characters of the string which are ", ". Otherwise we have a error in the query.
	multiple = multiple[:len(multiple)-2]
	// The database table needs to lose the last character "s" so we can use it to get the right table_id for the query.
	id := dbt[:len(dbt)-1]+"_id"
	// To math the amount of ? with the needed ID values, we create a new interface in which we append all the checked ID's,
	// as the first value and the first ?, we add the setTo variable which will tell the query what value it should use. Ex: trashed = 1 (move to trash) trashed = 0 (not in trash)
	// A interface can be read by the sql.Exec command as mutliple arguments. This is not possible from a Slice.
	args := []interface{}{setTo}
	for _,value := range checked {
		args = append(args, value)
		fmt.Println(args)

	}
	fmt.Println("UPDATE "+dbt+" SET "+row+" = ? WHERE "+id+" IN ("+multiple+")")
	// Prepare query with the right table name to update and the table_id, followed by the x amount of "?"
	// Ex: UPDATE posts SET trashed = 1 WHERE posts_id IN (?, ?, ?)
	query, err := db.Prepare("UPDATE "+dbt+" SET "+row+" = ? WHERE "+id+" IN ("+multiple+")")
	checkErr(err)
	// Execute query Ex: // Ex: UPDATE posts SET trashed = 1 WHERE post_id IN (1, 2, 3)
	_, err = query.Exec(args...)
	checkErr(err)

	http.Redirect(w, r, "/admin/"+dbt, http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, dbt string ){
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)
	// Getting the selected checkboxes from the request.Form
	// !!!! specify the checkbox- "files", ord "posts" etc, also in the views.
	// so if we have multiple sleections in different section, e.x admin page, then we wont have an error!
	checked := r.Form["checkbox"]
	placeholder,args := Multiple(checked)
	// The database table needs to lose the last character "s" so we can use it to get the right table_id for the query.

	if(dbt == "files"){
		DeleteFiles(w,r,placeholder,args)
	}

	id := dbt[:len(dbt)-1]+"_id"

	// Prepare query with the right table name to delete from and the table_id, followed by the x amount of "?"
	// Ex: DELete FROM posts WHERE post_id IN (?, ?, ?)
	query, err := db.Prepare("DELETE FROM "+dbt+" WHERE "+id+" IN ("+placeholder+")")
	checkErr(err)
	// Execute query Ex: // Ex: UPDATE posts SET trashed = 1 WHERE posts_id IN (1, 2, 3)
	_, err = query.Exec(args...)
	checkErr(err)

	http.Redirect(w, r, "/admin/"+dbt+"/trashed-"+dbt, http.StatusFound)
}

func DeleteFiles(w http.ResponseWriter, r *http.Request, placeholder string,values []interface {}) {
	var file_id int64
	var path string

	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	// add the placeholders
	rows, err := db.Query("SELECT file_id,path FROM files WHERE file_id IN("+placeholder+")",values...)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&file_id, &path)
		checkErr(err)
		// We use "/file/" as a prefix in the DB or else we get in conflict with the router which also use files,
		// when we serve static files the /file/ will be chaged to files.

		// if the folder path is fetched from the db it will state file/some/file/path
		// to remove the file we need to change the path to /files/ instead of file/
		// otherswise we use a incorrect filepath which will result in an error.

		// !! we do this path[5:]for the above reason.
		err = os.Remove("GolangBlog/static/files/"+path[5:])
		checkErr(err)
	}
}

func DeleteFolders(w http.ResponseWriter, r *http.Request, dbt string)(msg []string){
	//
	// DeleteFiles(w,r,dbt,args)
	var file_id int64
	var path string
	var folderPath string

	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	// add the placeholders
	checked := r.Form["checkbox"]
	fmt.Println("checked folders: ", checked)

	placeholder,values := Multiple(checked)
	rows, err := db.Query("SELECT file_id,path FROM "+dbt+" WHERE folder_id IN("+placeholder+")",values...)
	checkErr(err)
	// removing files from filesystem
	for rows.Next() {
		err = rows.Scan(&file_id, &path)
		checkErr(err)
		// We use "/file/" as a prefix in the DB or else we get in conflict with the router which also use files,
		// when we serve static files the /file/ will be chaged to files.

		// if the folder path is fetched from the db it will state file/some/file/path
		// to remove the file we need to change the path to /files/ instead of file/
		// otherswise we use a incorrect filepath which will result in an error.

		// !! we do this path[5:]for the above reason.
		err = os.Remove("GolangBlog/static/files/"+path[5:])
		checkErr(err)
	}
	//removing rows from database
	query, err := db.Prepare("DELETE FROM "+dbt+" WHERE folder_id IN ("+placeholder+")")
	checkErr(err)
	_, err = query.Exec(values...)
	checkErr(err)
	//removing folder from filesystem
	rows, err = db.Query("SELECT path FROM folders WHERE folder_id IN("+placeholder+")",values...)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&folderPath)
		checkErr(err)
		err := os.RemoveAll("GolangBlog/static/"+folderPath)
		checkErr(err)
		if (err == nil) {
			msg = append(msg,folderPath + "is removed successfully")
		} else {
			msg = append(msg,folderPath + "The folder you want to delete doesn't exist")
		}
	}
	//removing folder rows from database
	query, err = db.Prepare("DELETE FROM folders WHERE folder_id IN ("+placeholder+")")
	checkErr(err)
	_, err = query.Exec(values...)
	checkErr(err)

	//removeRows(checked[0])

	return msg
}

// Takes a slice of string ,[]string, which consist of the id's of the checked items in the form
// of course other []string can be provided to create prepared statements with multiple values.
func Multiple(multiple []string)(string,[]interface {}) {
	var placeholder string
	fmt.Println("multiple :", multiple)
	// Creating a string with the amount of ? required for the Query string by using the length of the checkbox slice.
	placeholder = strings.Repeat("?, ",len(multiple))
	// delete the last 2 characters of the string which are ", ". Otherwise we have a error in the query.
	placeholder = placeholder[:len(placeholder)-2]
	fmt.Println(placeholder)

	// the stmt.Exec can loop over an interface, so let's make one first to populate.
	values := make([]interface{}, 0, len(multiple))
	// To math the amount of ? with the needed ID values, we create a new interface in which we append all the checked ID's,
	// A interface can be read by the sql.Exec command as multiple arguments. This is not possible from a Slice.
	for _,v := range multiple {
		values = append(values, v)
	}
	return placeholder,values
}

// must be called inside the deleteFolders function after removing the dirs.
func removeRows(folderID string){
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT folder_id FROM folders WHERE parent_id = ?",folderID)
	checkErr(err)
	fmt.Println(err)
	// adding the folder to be deleted to the slice of string, below we will check for children
	// and if so, add them.
	folders := []string{folderID}
	// checks if the row from the db is not empty,
	// if not, selects the id to the parent id,row[id], so we can get,
	// all children from the top deleted album.
	/* Example:
	 * $id = 5 (Folder Users)
	 * $row[album_id] = 24 ( user admin has parent_id 5, the folder Users)
	 * Then again we check if there are folders with a parent_id of 24
	 * if there is, add it to the array of folder_id's to delete.
	 * In this case there is.
	 * $row['album_id'] = 22 (admins contacts folder) has a parent_id of 24
	 * This goes on until there are no folders left with a parent_id of 22 in this case.
	*/

	// wat als de parent folder meerdere children heeft? loopt hij dan vast?
	// slaat hij ze over? worden ze wel verwijderd?
	// DE FUNCTIE ZICHZELF LATEN aanroepen op het laatst, dan begint het steeds opnieuw.
	// if err != nil {if err != nil {
	for err == nil {
		for rows.Next() {
			err = rows.Scan(&folderID)
			folders = append(folders, folderID)
		}
		rows, err = db.Query("SELECT folder_id FROM folders WHERE parent_id = ?", folderID);
	}
	fmt.Println("Folders slice: ", folders)
	// because we now have a new row[folderID], we need to check again if its empty,
	// if it is not, push it to the array.
	//if it is, don't push it, en the loop will end with the while clause.
	//if(!empty($row['album_id'])){
	//$folders_id[] = $row['album_id'];
	//}

	// Create s string with all the album id's.
	// deleting rows from database.
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}