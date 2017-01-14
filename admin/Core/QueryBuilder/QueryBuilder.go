package QueryBuilder

import (
	"strings"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"log"
	"fmt"
)

type Query struct{
	PrimaryKey string
	Table string
	Query string
	Allowed []string
	Relations []string
	Values []string
	Columns []string
}

func (m *Query) All() {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	defer db.Close()
	query := m.Select()+m.From()+m.Where()+m.OrderBy();

	fmt.Println("Query: ", query)
}

func (m *Query) Select() string {
	return "SELECT * "
}
func (m *Query) From() string{
	return "FROM "+m.Table
}
func (m *Query) Where() string {
	parts := []string{}
	columns := map[string]string{ m.PrimaryKey : "10" }
	for column,value := range columns {
		// Set column to ? for prepared statement
 		parts = append(parts,m.Table+"."+column+" = ?")
		// Set values array for PDO prepared statement
		m.Values = append(m.Values,value);
	}
	where := strings.Join(parts," AND ")
	return " WHERE "+where;
}

func (m *Query) OrderBy() string{
	return " ORDER BY "+m.Table+"."+m.PrimaryKey+" DESC"
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}