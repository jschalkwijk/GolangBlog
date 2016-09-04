package files

import (
	"database/sql"
	_"database/sql/driver"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"net/http"
	"fmt"
	"strings"
	"os"
)


func Delete(w http.ResponseWriter, r *http.Request, placeholder string,values []interface {}) {
	var file_id int64
	var path string

	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	// add the placeholders
	rows, err := db.Query("SELECT file_id,path FROM files WHERE file_id IN("+placeholder+")",values)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&file_id, &path)
		checkErr(err)
		os.Remove(path)
	}
}

func multiple(multiple []string)(string,[]interface {}) {
	var placeholder string
	// the values are stored in the r.Form["checked files"], it's a slice of strings.
	// for the lengths of the slice we need to create the ?, for the prepared statement.
	placeholder = strings.Repeat("?, ",len(multiple))
	// When the placeholders string is populated we need to remove the trailing (, ) including the space.
	// removing the last 2 characters of the placeholder by checking it's length minus 2.
	placeholder = placeholder[:len(placeholder)-2]
	fmt.Println(placeholder)

	// the stmt.Exec can loop over an interface, so let's make one first to populate.
	values := make([]interface{}, 0, len(multiple))
	// Add every value of the checked-files []string to the values interface
	for _,v := range multiple{
		values = append(values, v)
	}
	return placeholder,values
}

