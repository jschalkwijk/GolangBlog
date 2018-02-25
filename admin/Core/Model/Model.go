package Model

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"log"
	"github.com/jmoiron/sqlx"
	"fmt"
)

/*
	TODO:
	Een interface die refreeerd naar alle common functions die elke model moet hebben
	De hoof model functie neemt dan elke model aan,het model dat doorgegeven wordt
	moet de table name hebben etc, zo kan ik mischien general functies maken voor
	crud functionaliteiten.
*/


type BaseModel interface {
	Execute(query string) (*sqlx.Rows,error)
}

type Model struct{
	ID string
	PrimaryKey string
	Table string
	Allowed map[string]int
	Relations []string
}
func (m *Model) Execute(query string, values []interface{}) (*sqlx.Rows,error){
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()
	var rows *sqlx.Rows

	if len(values) < 1 {
		// Selects all rows from roles, and links the category_id row to the matching title.
		rows, err = db.Queryx(query)
	} else {
		rows, err = db.Queryx(query,values...)
	}

	checkErr(err)

	return rows,err
}
func (m *Model) PrepareExecute(query string, values []interface{}) (error){
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare(query)
	fmt.Println(stmt)
	checkErr(err)
	/* To be able to save the new html to the database, convert it to a slice of bytes, why is this working?, we can't save
	 * a value of type template.HTML to the DB. I tried different things, change the .Content to string, byte, but then I have a problem displaying
	 * the content in html format on the page.
	 */
	_, err = stmt.Exec(values...)
	checkErr(err)

	return err
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}