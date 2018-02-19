package QueryBuilder

import (
	"strings"
	_"github.com/go-sql-driver/mysql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"log"
	"go/types"
	"fmt"
	"github.com/jmoiron/sqlx"
)

/*
	TODO:
	Een interface die refreeerd naar alle common functions die elke model moet hebben
	De hoof model functie neemt dan elke model aan,het model dat doorgegeven wordt
	moet de table name hebben etc, zo kan ik mischien general functies maken voor
	crud functionaliteiten.
*/
type Sql interface {
	All()
	Select(columns types.Slice) string
	From() string
	Where() string

}
type Query struct{
	ID string
	PrimaryKey string
	Table string
	Query string
	Allowed []string
	Relations []string
	Values []string
	Columns []string
}
func (q *Query) Execute(query string) (*sqlx.Rows,error){
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	// Selects all rows from roles, and links the category_id row to the matching title.
	rows, err := db.Queryx(query)
	checkErr(err)

	return rows,err
}
func (q *Query) All() (*sqlx.Rows,error) {
	query := q.Select([]string{})+q.From()+q.OrderBy();
	fmt.Println(query)
	return q.Execute(query)
}

func (q *Query) One() *sqlx.Row {
	//query := q.Select([]string{})+q.From()+q.Where()+q.OrderBy();
	query := "SELECT * FROM "+q.Table+" WHERE "+q.PrimaryKey+" = ? LIMIT  1"
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	row := db.QueryRowx(query, q.ID)

	fmt.Println(query,q.ID)
	return row
}

func (q *Query) Select(columns []string ) string {
	if len(columns) > 0 {
		return "SELECT columns "
	} else {
		return "SELECT * "
	}

}
func (q *Query) From() string{
	return "FROM "+q.Table
}
func (q *Query) Where() string {
	parts := []string{}
	columns := map[string]string{ q.PrimaryKey : "10" }
	for column,value := range columns {
		// Set column to ? for prepared statement
 		parts = append(parts,q.Table+"."+column+" = ?")
		// Set values array for PDO prepared statement
		q.Values = append(q.Values,value);
	}
	where := strings.Join(parts," AND ")
	return " WHERE "+where;
}

func (q *Query) OrderBy() string{
	return " ORDER BY "+q.Table+"."+q.PrimaryKey+" DESC"
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}