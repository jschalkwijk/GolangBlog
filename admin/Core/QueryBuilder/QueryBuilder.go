package QueryBuilder

import (
	"fmt"
	"strings"
	"log"
	"github.com/jschalkwijk/GolangBlog/admin/Core/Model"
	"github.com/jmoiron/sqlx"
)

type Query struct {
	Model.Model
	Values []interface{}
	Columns []string
}

func (q *Query) All() (*sqlx.Rows,error) {
	query := q.Select([]string{})+q.From(q.Table)+q.OrderBy();
	fmt.Println(query)
	return q.Execute(query,[]interface{}{})
}

func (q *Query) One(id string,model interface{}) (interface{}) {
	//query := m.Select([]string{})+m.From()+m.Where()+m.OrderBy();

	query := q.Select([]string{})+q.From(q.Table)+q.Where(map[string]string{q.PrimaryKey:q.ID})+" LIMIT 1"


	fmt.Println(query)
	fmt.Println(q.Columns)
	fmt.Println(q.Values)

	rows,err := q.Execute(query,q.Values)

	for rows.Next() {
		err = rows.StructScan(
			model,
		)
		checkErr(err)
	}
	return &model
}

func (q *Query) Select(columns []string ) string {
	if len(columns) > 0 {
		return "SELECT columns "
	} else {
		return "SELECT * "
	}

}
func (q *Query) From(table string) string{
	q.Table = table
	return "FROM "+table
}
func (q *Query) Where(columns map[string]string) string {
	parts := []string{}
	fmt.Println(columns)
	q.Values = make([]interface{}, 0, len(columns))
	for column,value := range columns {

		q.Values = append(q.Values, value);
		// Set Columns Slice
		q.Columns = append(q.Columns,column)
		// Set column to ? for prepared statement
		parts = append(parts,q.Table+"."+column+" = ?")

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