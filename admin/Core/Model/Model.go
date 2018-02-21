package Model

import (
	_"github.com/go-sql-driver/mysql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"log"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jschalkwijk/GolangBlog/admin/Core/QueryBuilder"
)

/*
	TODO:
	Een interface die refreeerd naar alle common functions die elke model moet hebben
	De hoof model functie neemt dan elke model aan,het model dat doorgegeven wordt
	moet de table name hebben etc, zo kan ik mischien general functies maken voor
	crud functionaliteiten.
*/

type BaseModel interface {
	All() (*sqlx.Rows,error)
	One(id string) *sqlx.Row
}

type Model struct{
	QueryBuilder.Query
	ID string
	PrimaryKey string
	Table string
	Allowed []string
	Relations []string
}
func (m *Model) Execute(query string) (*sqlx.Rows,error){
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	// Selects all rows from roles, and links the category_id row to the matching title.
	rows, err := db.Queryx(query)
	checkErr(err)

	return rows,err
}
func (m *Model) All() (*sqlx.Rows,error) {
	var q = m.Query

	query := q.Select([]string{})+q.From(m.Table)+q.OrderBy();
	fmt.Println(query)
	return m.Execute(query)
}

func (m *Model) One(id string) *sqlx.Row {
	//query := m.Select([]string{})+m.From()+m.Where()+m.OrderBy();

	query := m.Select([]string{})+m.From(m.Table)+m.Where(map[string]string{m.PrimaryKey:m.ID,"description":"Jorn"})+m.OrderBy()
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	row := db.QueryRowx(query, m.Values...)

	fmt.Println(query)
	fmt.Println(m.Columns)
	fmt.Println(m.Values)

	return row
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}