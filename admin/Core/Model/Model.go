package Model

import (
	"strings"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"reflect"
	"fmt"
	"log"
)

type Model struct{
	PrimaryKey string
	Table string
	Query string
	Allowed []string
	Relations []string
	Values []string
	Columns []string
}

var typeRegistry = make(map[string]reflect.Type)

func (m *Model) All() *sql.Rows {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	defer db.Close()
	query := m.Select()+m.From()+m.Where()+m.OrderBy();

	println("Query: ", query)
	println("Values: ", m.Values[0])
	// Selects all rows from posts, and links the category_id row to the matching title.
	rows, err := db.Query(query,m.Values[0])
	checkErr(err)
	//columns, err := rows.Columns()
	//for key,column := range columns{
	//	println(key,column)
	//	m.Columns = append(m.Columns,column)
	//}

	return rows

	//// columns = []string{"Post_ID","Title"}
	//fieldindex := []int{}
	//for i := 0; i < len(m.Columns); i++{
	//	typeRegistry[m.Columns[i]] = reflect.TypeOf(m.Columns[i])
	//	fieldindex = append(fieldindex,i)
	//}
	//fmt.Printf("%s",typeRegistry)
	//
	//for _,index := range fieldindex {
	//	v = reflect.Indirect(v).Field(index)
	//	//deze wordt niet uitgevoerd
	//
	//	alloc := reflect.New(typeRegistry[m.Columns[index]])
	//	fmt.Printf("%s", reflect.TypeOf(m.Columns[index]))
	//	fmt.Printf("%s", alloc.Kind())
	//	fmt.Println(reflect.ValueOf(alloc))
	//	//v.Set(alloc)
	//
	//
	//	if v.Kind() == reflect.Map && v.IsNil() {
	//		v.Set(reflect.MakeMap(typeRegistry[m.Columns[index]]))
	//	}
	//}
	//
	//for i := 0; i < v.NumField(); i++ {
	//	f := v.Field(i)
	//	fmt.Printf("%d: %s %s = %v\n",i,TypeOfModel.Field(i).Name,f.Type(),f.Interface())
	//}
}

func (m *Model) Rape(model interface{}){
	// the stmt.Exec can loop over an interface, so let's make one first to populate.

	values := make([]interface{},12)
	// To math the amount of ? with the needed ID values, we create a new interface in which we append all the checked ID's,
	// A interface can be read by the sql.Exec command as multiple arguments. This is not possible from a Slice.
	//for _,v := range columns {
	//	values = append(values, v)
	//}
	v := reflect.ValueOf(model)
	//v is currently a pointer to *model.Model
	// indirect return the struct v is pointing to.
	// to get the struct fields key,type,values we need  struct not the pointer
	v = reflect.Indirect(v).Elem()
	TypeOfModel := v.Type()
	fmt.Println(TypeOfModel)
	for i := 1; i < v.NumField(); i++ {
		f := v.Field(i)
		fmt.Printf("%d: %s %s = %v\n",i,TypeOfModel.Field(i).Name,f.Type(),f.Interface())
		values[i] = v.Field(i).Addr().Interface()
	}
	fmt.Printf("%s",values)
	rows := m.All()
	for rows.Next() {
		err := rows.Scan(values...)
		checkErr(err)
	}

	v = reflect.ValueOf(model)
	//v is currently a pointer to *model.Model
	// indirect return the struct v is pointing to.
	// to get the struct fields key,type,values we need  struct not the pointer
	v = reflect.Indirect(v).Elem()
	TypeOfModel = v.Type()
	fmt.Println(TypeOfModel)
	for i := 1; i < v.NumField(); i++ {
		f := v.Field(i)
		fmt.Printf("%d: %s %s = %v\n",i,TypeOfModel.Field(i).Name,f.Type(),f.Interface())
	}

}
func (m *Model) Select() string {
	return "SELECT * "
}
func (m *Model) From() string{
	return "FROM "+m.Table
}
func (m *Model) Where() string {
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

func (m *Model) OrderBy() string{
	return " ORDER BY "+m.Table+"."+m.PrimaryKey+" DESC"
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

//func StructScan(dest interface{}) error {
//	v := reflect.ValueOf(dest)
//
//	if v.Kind() != reflect.Ptr {
//		return errors.New("must pass a pointer, not a value, to StructScan destination")
//	}
//
//	v = reflect.Indirect(v)
//
//	if !r.started {
//		columns, err := r.Columns()
//		if err != nil {
//			return err
//		}
//		m := r.Mapper
//
//		r.fields = m.TraversalsByName(v.Type(), columns)
//		// if we are not unsafe and are missing fields, return an error
//		if f, err := missingFields(r.fields); err != nil && !r.unsafe {
//			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
//		}
//		r.values = make([]interface{}, len(columns))
//		r.started = true
//	}
//
//	err := fieldsByTraversal(v, r.fields, r.values, true)
//	if err != nil {
//		return err
//	}
//	// scan into the struct field pointers and append to our results
//	err = r.Scan(r.values...)
//	if err != nil {
//		return err
//	}
//	return r.Err()
//}
//
//// Connect to a database and verify with a ping.
//func Connect(driverName, dataSourceName string) (*DB, error) {
//	db, err := Open(driverName, dataSourceName)
//	if err != nil {
//		return db, err
//	}
//	err = db.Ping()
//	return db, err
//}