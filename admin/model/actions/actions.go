package actions

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	cfg "github.com/jschalkwijk/GolangBlog/admin/config"
	"fmt"
	"strings"
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
	db, err := sql.Open("mysql", cfg.DB)
	defer db.Close()
	checkErr(err)
	fmt.Println("starts")
	// Getting the selected checkboxes from the request.Form
	checked := r.Form["checkbox"]
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
	db, err := sql.Open("mysql", cfg.DB)
	defer db.Close()
	checkErr(err)
	// Getting the selected checkboxes from the request.Form
	checked := r.Form["checkbox"]
	// Creating a string with the amount of ? required for the Query string by using the length of the checkbox slice.
	multiple := strings.Repeat("?, ",len(checked))
	// delete the last 2 characters of the string which are ", ". Otherwise we have a error in the query.
	multiple = multiple[:len(multiple)-2]
	// The database table needs to lose the last character "s" so we can use it to get the right table_id for the query.
	id := dbt[:len(dbt)-1]+"_id"
	// To math the amount of ? with the needed ID values, we create a new interface in which we append all the checked ID's,
	// A interface can be read by the sql.Exec command as mutliple arguments. This is not possible from a Slice.
	args := []interface{}{}
	for _,value := range checked {
		args = append(args, value)
		fmt.Println(args)

	}
	fmt.Println("DELETE FROM "+dbt+" WHERE "+id+" IN ("+multiple+")")
	// Prepare query with the right table name to delete from and the table_id, followed by the x amount of "?"
	// Ex: DELete FROM posts WHERE post_id IN (?, ?, ?)
	query, err := db.Prepare("DELETE FROM "+dbt+" WHERE "+id+" IN ("+multiple+")")
	checkErr(err)
	// Execute query Ex: // Ex: UPDATE posts SET trashed = 1 WHERE posts_id IN (1, 2, 3)
	_, err = query.Exec(args...)
	checkErr(err)

	http.Redirect(w, r, "/admin/"+dbt+"/trash", http.StatusFound)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}