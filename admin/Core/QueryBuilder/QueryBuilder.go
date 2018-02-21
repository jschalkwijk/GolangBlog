package QueryBuilder

import (
	"fmt"
	"strings"
)

type Query struct {
	ID string
	PrimaryKey string
	Table string
	Allowed []string
	Relations []string
	Values []interface{}
	Columns []string
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