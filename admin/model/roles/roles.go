package roles

import (
	_"github.com/go-sql-driver/mysql"
	//"database/sql"
	"fmt"
	//"html/template"
	//"net/http"
	//"strconv"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"github.com/jschalkwijk/GolangBlog/admin/Core/QueryBuilder"
	//"github.com/gorilla/schema"
	//"github.com/jmoiron/sqlx"
	"log"
	"database/sql"
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

var Query QueryBuilder.Query

func init(){
	Query.PrimaryKey = "role_id"
	Query.Table = "roles"
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

	rows,err := Query.All()
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
	Query.ID = id

	rows := Query.One()
	data := new(Data)
	role := new(Role)

	err := rows.StructScan(
		&role,
	)
	checkErr(err)
	// convert string to HTML markdown

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
func (r *Role) update() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)

	fmt.Println("reference to Role struct: ", r)

	stmt, err := db.Prepare("UPDATE roles SET name=?, description=?, WHERE role_id=?")
	fmt.Println(stmt)
	checkErr(err)
	/* To be able to save the new html to the database, convert it to a slice of bytes, why is this working?, we can't save
	 * a value of type template.HTML to the DB. I tried different things, change the .Content to string, byte, but then I have a problem displaying
	 * the content in html format on the page.
	 */
	_, err = stmt.Exec(r.Name,r.Description,r.Role_ID)
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
