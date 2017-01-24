/*	-- Categories Model --
 * 	All functions in this file are called by the corresponding controller or by
 	functions from itself.
 */
package categories

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"
	cfg "github.com/jschalkwijk/GolangBlog/admin/config"
	"log"
	"github.com/gorilla/schema"
)

/* Category struct will hold data about a category and can be added to the Data struct */
type Category struct {
	Category_ID int `schema:"-"`
	Title string
	Description string
	Content string
	Keywords string
	Approved int
	Author string
	Cat_Type string
	Date string
	Parent_ID int
	Trashed int
}

/*
 * Declaring vars corresponding to the struct. When scanning data from the database, the
   data will be stored on the memory address of these vars.
*/

/* Stores a single category, or multiple categories in a Slice which can be iterated over in the template */
type Data struct {
	Categories []*Category
	Message string
	Deleted bool
}

/* -- Get all categories --
 * 	Connects to the database and gets all category rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	For every row get the values and set the values to the memory address of the named variable.
 		- Instantiate a new Category Struct and insert values.
 		- Append the category struct to the Data.Categories Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func All(trashed int) *Data {
	db, err := sql.Open("mysql", cfg.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	rows, err := db.Query("SELECT * FROM categories WHERE trashed = ? ORDER BY categorie_id DESC", trashed)
	checkErr(err)

	data := new(Data)

	for rows.Next() {
		c := new(Category)
		err = rows.Scan(
			&c.Category_ID,
			&c.Title,
			&c.Description,
			&c.Content,
			&c.Keywords,
			&c.Approved,
			&c.Author,
			&c.Cat_Type,
			&c.Date,
			&c.Parent_ID,
			&c.Trashed,
		)
		checkErr(err)
		data.Categories = append(data.Categories , c)
	}

	if(trashed == 1) {
		data.Deleted = true
	} else {
		data.Deleted = false
	}

	return data
}

/* -- Get a single categories -- */
/* GetOneCategory gets a category from the DB and returns a pointer to the Struct. It takes a id and category_title.
 * 	Connects to the database and gets all category rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new Category Struct and insert values.
 *	Append the category struct to the Data.Categories Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func One(id string) *Data {
	db, err := sql.Open("mysql", cfg.DB)
	checkErr(err)

	defer db.Close()

	rows := db.QueryRow("SELECT * FROM categories WHERE categorie_id=? LIMIT 1", id)

	data := new(Data)
	c := new(Category)

	err = rows.Scan(
		&c.Category_ID,
		&c.Title,
		&c.Description,
		&c.Content,
		&c.Keywords,
		&c.Approved,
		&c.Author,
		&c.Cat_Type,
		&c.Date,
		&c.Parent_ID,
		&c.Trashed,
	)

	checkErr(err)

	data.Categories = append(data.Categories , c)

	//fmt.Println(collection.categories)
	return data
}

/* -- Category Methods -- */

/* saveCategory updates the values of an existing category to the database and is a method to Category
 * Called by EditCategory
 * Connect to the DB and prepares query.
 * Execute query with the inserted struct values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (c *Category) update() error {
	db, err := sql.Open("mysql", cfg.DB)
	checkErr(err)

	defer db.Close()

	stmt, err := db.Prepare("UPDATE categories SET title=?, description=? WHERE categorie_id=?")
	fmt.Println(stmt)
	checkErr(err)
	_, err = stmt.Exec(c.Title,c.Description,c.Category_ID)
	checkErr(err)

	return err
}

/* AddCategory saves the values of a new category to the database and is a method to Category.
 * Called by CreateCategory
 * Connect to the DB and prepares query.
 * Execute query with the inserted values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (c *Category) Store() error {
	db, err := sql.Open("mysql", cfg.DB)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO categories (title,description) VALUES(?,?) ")
	checkErr(err)
	_, err = stmt.Exec(c.Title,c.Description)

	checkErr(err)
	return err
}
// End category methods

/* EditCategory takes updated form values from the http.request to populate a Category and call the saveCategory method.
 * The request delivers the FormValues if asked.
 * Convert category_id to an INT. The category ID is pulled from the from as a string.
 * FormValues are appointed to to the memory address of the Category struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call saveCategory, a method of the Category Struct, to update the DB
*/
func (c *Category) Patch(r *http.Request) (*Data,bool) {
	data:= new(Data)
	updated := false

	if r.Method == "POST" {
		err := r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(c,r.PostForm)
		checkErr(err)
		if c.Title == "" {
			data.Message = "Please fill in all the required fields"
		} else {
			err = c.update()
			checkErr(err)
			updated = true
		}
	}

	data.Categories = append(data.Categories,c)

	return data,updated
}

/* NewCategory takes updated form values from the http.request to populate a Category and call the AddCategory method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the Category struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call AddCategory, a method of the Category Struct, to insert new category in the DB.
*/
func Create(r *http.Request) (*Data,bool) {
	c := new(Category)
	data := new(Data)
	created := false

	if r.Method == "POST"{
		err := r.ParseForm()
		checkErr(err)
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(c,r.PostForm)
		checkErr(err)
		if c.Title == "" {
			data.Message = "Please fill in all the required fields"
		} else {
			err = c.Store()
			checkErr(err)
			created = true
		}
		data.Categories = append(data.Categories , c)

	} else{
		data.Categories = append(data.Categories,c)
	}

	return data,created
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

