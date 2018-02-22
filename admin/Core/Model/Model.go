package Model

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"log"
	"github.com/jmoiron/sqlx"
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
	Allowed []string
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

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}