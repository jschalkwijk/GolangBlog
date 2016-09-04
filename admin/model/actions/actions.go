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

// Takes a slice of string ,[]string, which consist of the id's of the checked items in the form
// of course other []string can be provided to create prepared statements with multiple values.
func Multiple(multiple []string)(string,[]interface {}) {
	var placeholder string
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}