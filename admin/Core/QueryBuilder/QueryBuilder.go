package QueryBuilder

import (
	"fmt"
	"strings"
	"log"
	"github.com/jschalkwijk/GolangBlog/admin/Core/Model"
	//"reflect"
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

//func (q *Query) All(model interface{}) ([]interface{}){
//	query := q.Select([]string{})+q.From(q.Table)+q.OrderBy();
//	fmt.Println(query)
//	rows,err := q.Execute(query,[]interface{}{})
//	checkErr(err)
//	models := make([]interface{},0,0)
//	slice := reflect.ValueOf(model).Elem()
//	value := reflect.ValueOf(model)
//	// Allocate slice with desired capacity
//	fmt.Printf("Hello, Model: %s\n", slice)
//	fmt.Print(reflect.TypeOf(model))
//
//	for rows.Next() {
//		role := reflect.New(value.Type())
//		role.Elem().Set(value)
//
//		//role := reflect.ValueOf(value).Elem()
//		fmt.Printf("Hello, item2 : %s\n", role.String())
//		//fmt.Println(value)
//		// puts all columns inside the Role struct automaticly
//		err := rows.StructScan(
//			role.Interface(),
//		)
//		checkErr(err)
//		//fmt.Println(role.Role_ID,role.Title)
//		fmt.Println(role)
//		// Add the Role to the Roles slice.
//		models = append(models , role)
//	}
//	return models
//}

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
	return model
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