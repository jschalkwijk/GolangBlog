package roles

import (
	_"github.com/go-sql-driver/mysql"
	//"database/sql"
	"fmt"
	//"html/template"
	//"net/http"
	//"strconv"

	//"github.com/gorilla/schema"
	//"github.com/jmoiron/sqlx"
	"log"
	"github.com/jschalkwijk/GolangBlog/admin/Core/QueryBuilder"
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
)
/* Role struct will hold data about a role and can be added to the Data struct */
type Role struct {
	Role_ID int `schema:"-"`
	Name string
	Description string
	Created_At string `schema:"-"`
	Updated_At string `schema:"-"`
}

/* Stores a single role, or multiple roles in a Slice which can be iterated over in the template */
type Data struct {
	Roles      []*Role
	Message string
}

var q QueryBuilder.Query

func init(){
	// we don't need to assign the embedded model struct but for readability I chose to do so.
	// Also it now point to the memory address.
	m := &q.Model
	m.PrimaryKey = "role_id"
	m.Table = "roles"
	m.Allowed = map[string]int{"Name":0,"Description":1}
}
/* -- Get all Roles --
 * 	Connects to the database and gets all roles rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	For every row get the values, and set the values to the memory address of the named variable.
 		- Instantiate a new Role Struct and insert values.
 		- Append the Role struct to the Data.Roles Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func All() *Data {

	rows,err := q.All()
	checkErr(err)

	data := new(Data)

	for rows.Next() {
		role := new(Role)
		// puts all columns inside the Role struct automaticly
		err := rows.StructScan(
			&role,
		)
		checkErr(err)
		//fmt.Println(role.Role_ID,role.Title)
		fmt.Println(role)
		// Add the Role to the Roles slice.
		data.Roles = append(data.Roles , role)
	}

	return data
}
//func All() *Data {
//
//	rows := q.All(new(Role))
//	fmt.Printf("%s",rows)
//	data := new(Data)
//
//	for _, role := range rows {
//		fmt.Println(role)
//		// Add the Role to the Roles slice .
//		fmt.Println(reflect.TypeOf(role))
//		data.Roles = append(data.Roles ,role)
//	}
//
//	return data
//}

/* -- Get a single Role -- */
/* GetOneRole gets a role from the DB and returns a pointer to the Struct. It takes a id and role_title.
 * 	Connects to the database and gets all role rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new Role Struct and insert values.
 *	Append the Role struct to the Data.Roles Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func One(id string, getPermission bool) *Data {
	q.ID = id

	model := new(Role)
	role := q.One(id,model).(*Role)

	data := new(Data)

	data.Roles = append(data.Roles , role)


	/* When we need to edit or create a role, we need to get the permissions in order to select them inside the html page.
	 * since we already have a function inside the permissions model, we will call that.
	 * This returns a pointer to the Data struct of model/roles. We  set our
	   roles.Data struct, Roles, to the slice of roles.Data.Permissions.
	 * They are accessible inside the template now.
	*/
	//if(getPermission) {
	//	listCat := cat.All(0)
	//	data.Categories = listCat.Categories
	//}
	//
	//fmt.Println(data.Categories)
	return data
}

/* -- Role Methods -- */

/* update updates the values of an existing role to the database and is a method to Role
 * Called by EditRole
 * Connect to the DB and prepares query.
 * Execute query with the inserted struct values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (role *Role) Patch(r *http.Request) (*Data,bool) {
	data := new(Data)
	updated := false

	err := r.ParseForm()

	decoder := schema.NewDecoder()
	decoder.ZeroEmpty(true)
	err = decoder.Decode(role, r.PostForm)
	checkErr(err)
	columns := r.Form
	fmt.Println(columns)
	//columns := map[string]string{
	//	"name":role.Name,
	//	"description":role.Description,
	//	"role_id": role.Role_ID,
	//}
	err = role.Update(columns)
	checkErr(err)
	data.Roles = append(data.Roles , role)

	return data,updated
}

func (r *Role) Update(role url.Values) error {
		err := q.Update(role)
		checkErr(err)
		return err
}
/* save saves the values of a new category to the database and is a method to Role.
 * Called by CreateRole
 * Connect to the DB and prepares query.
 * Execute query with the inserted values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
//func (p *Role) save() error {
//	db, err := sql.Open("mysql", config.DB)
//	defer db.Close()
//	stmt, err := db.Prepare("INSERT INTO roles (title,description,content,category_id) VALUES(?,?,?,?)")
//	fmt.Println(stmt)
//	checkErr(err)
//	res, err := stmt.Exec(p.Title,p.Description,[]byte(p.Content),p.Category_ID)
//	fmt.Println(res)
//	checkErr(err)
//	return err
//}
// End Role methods

/* EditRole takes updated form values from the http.request to populate a Role and call the saveRole method.
 * The request delivers the FormValues if asked.
 * Convert role_id to an INT. The role ID is pulled from the from as a string.
 * FormValues are appointed to to the memory address of the Role struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call saveRole, a method of the Role Struct, to update the DB
*/
//func (role *Role) Patch(r *http.Request) (*Data,bool) {
//	data := new(Data)
//	updated := false
//
//	if r.Method == "POST" {
//		category_id := r.FormValue("Category_ID")
//		category := r.FormValue("Category")
//		/* 	To add a new category from a add role form we need to create a new
//			 category, and then get the new ID of that category to insert it into the Role struct.
//			 Also see addCategoryFromForm
//		 */
//		if (category != "") {
//			category_id = addCategoryFromForm(category, category_id);
//		} else {
//			fmt.Println("empty string")
//		}
//		//// convert string values to INT before inserting into Struct and DB
//		categoryINT, _ := strconv.Atoi(category_id)
//		err := r.ParseForm()
//		decoder := schema.NewDecoder()
//		decoder.ZeroEmpty(true)
//		err = decoder.Decode(role, r.RoleForm)
//
//		checkErr(err)
//
//		role.Category_ID = categoryINT;
//		if role.Title == "" || role.Content == "" {
//			fmt.Println("I got here! ")
//			data.Roles = append(data.Roles , role)
//			data.Message = "Please fill in all the required fields"
//			fmt.Println(data.Message)
//		} else {
//			err = role.update()
//			checkErr(err)
//			data.Roles = append(data.Roles , role)
//			updated = true
//		}
//	} else {
//		data.Roles = append(data.Roles , role)
//	}
//	return data,updated
//}

/* NewRole takes updated form values from the http.request to populate a Role and call the addRole method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the Role struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call addRole, a method of the Role Struct, to insert new role in the DB.
*/
//func Create(r *http.Request) (*Data,bool){
//
//	data := new(Data)
//	role := new(Role)
//
//	created := false
//
//	if r.Method == "POST" {
//		category_id := r.FormValue("Category_ID")
//		category := r.FormValue("Category")
//		/* 	To add a new category from a add role form we need to create a new
//			 category, and then get the new ID of that category to insert it into the Role struct.
//			 Also see addCategoryFromForm
//		 */
//		if (category != "") {
//			category_id = addCategoryFromForm(category, category_id);
//		} else {
//			fmt.Println("empty string")
//		}
//		//// convert string values to INT before inserting into Struct and DB
//		categoryINT, _ := strconv.Atoi(category_id)
//		err := r.ParseForm()
//		decoder := schema.NewDecoder()
//		decoder.ZeroEmpty(true)
//		err = decoder.Decode(role, r.RoleForm)
//		checkErr(err)
//
//		role.Category_ID = categoryINT;
//		fmt.Printf("%v",role.Title)
//		if role.Title == "" || role.Content == "" {
//			data.Roles = append(data.Roles , role)
//			data.Message = "Please fill in all the required fields"
//		} else {
//			err = role.save()
//			checkErr(err)
//			data.Roles = append(data.Roles , role)
//			created = true
//		}
//	} else {
//		data.Roles = append(data.Roles , role)
//	}
//	return data,created
//}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
