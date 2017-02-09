/*	-- Posts Model --
 * 	All functions in this file are called by the corresponding controller or by
 	functions from itself.
 */
package posts

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	cat "github.com/jschalkwijk/GolangBlog/admin/model/categories"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"log"
)

//func init(){
//	Model.GetAll(&Post{Post_ID:1,Title:"Reflection Test"})
//}

/* Post struct will hold data about a post and can be added to the Data struct */
type Post struct {
	Post_ID int `schema:"-"`
	Title string
	Description string
	Content template.HTML
	Keywords string `schema:"-"`
	Approved int `schema:"-"`
	Author string `schema:"-"`
	Date string `schema:"-"`
	Category_ID int
	Category string
	Trashed int `schema:"-"`
}

/* Stores a single post, or multiple posts in a Slice which can be iterated over in the template */
type Data struct {
	Posts      []*Post
	Categories []*cat.Category
	Deleted	bool
	Dashboard bool
	Message string
}

/* -- Get all Posts --
 * 	Connects to the database and gets all posts rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	For every row get the values, and set the values to the memory address of the named variable.
 		- Instantiate a new Post Struct and insert values.
 		- Append the Post struct to the Data.Posts Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func All(trashed int) *Data {
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	// Selects all rows from posts, and links the category_id row to the matching title.
	rows, err := db.Queryx("SELECT posts.*, categories.title AS category FROM categories JOIN posts ON categories.category_id = posts.category_id WHERE posts.trashed = ? ORDER BY posts.post_id DESC",trashed)
	checkErr(err)

	data := new(Data)

	var content string
	for rows.Next() {
		post := new(Post)
		//post.Table = "posts"
		//post.PrimaryKey = "post_id"
		err = rows.StructScan(
			&post,
		)
		checkErr(err)
		// convert string to HTML markdown
		post.Content = template.HTML(content)
		//fmt.Println(post.Post_ID,post.Title)
		fmt.Println(post)
		data.Posts = append(data.Posts , post)
		data.Dashboard = false
	}
	//for post := range data.Posts {
	//	fmt.Println(post)
	//}
	if(trashed == 1) {
		data.Deleted = true
	} else {
		data.Deleted = false
	}

	return data
}

/* -- Get a single Post -- */
/* GetOnePost gets a post from the DB and returns a pointer to the Struct. It takes a id and post_title.
 * 	Connects to the database and gets all post rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new Post Struct and insert values.
 *	Append the Post struct to the Data.Posts Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func One(id string, getCat bool) *Data {
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	rows := db.QueryRowx("SELECT posts.*, categories.title AS category FROM categories JOIN posts ON categories.category_id = posts.category_id WHERE post_id=? LIMIT  1", id)

	data := new(Data)
	post := new(Post)
	var content string

	err = rows.StructScan(
		&post,
	)
	checkErr(err)
	// convert string to HTML markdown
	post.Content = template.HTML(content)
	data.Posts = append(data.Posts , post)
	 /* When we need to edit or create a post, we need to get the categories in order to select them inside the html page.
	  * since we already have a function inside the categories model, we will call that.
	  * This returns a pointer to the Data struct of model/categories. We  set our
	    categories.Data struct, Categories, to the slice of posts.Data.Categories.
	  * They are accessible inside the template now.
	 */
	if(getCat) {
		listCat := cat.All(0)
		data.Categories = listCat.Categories
	}

	fmt.Println(data.Categories)
	return data
}

/* -- Post Methods -- */

/* savePost updates the values of an existing post to the database and is a method to Post
 * Called by EditPost
 * Connect to the DB and prepares query.
 * Execute query with the inserted struct values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (p *Post) update() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)

	fmt.Println("reference to Post struct: ", p)

	stmt, err := db.Prepare("UPDATE posts SET title=?, description=?,category_id=?, content=? WHERE post_id=?")
	fmt.Println(stmt)
	checkErr(err)
	/* To be able to save the new html to the database, convert it to a slice of bytes, why is this working?, we can't save
	 * a value of type template.HTML to the DB. I tried different things, change the .Content to string, byte, but then I have a problem displaying
	 * the content in html format on the page.
	 */
	_, err = stmt.Exec(p.Title,p.Description,p.Category_ID,[]byte(p.Content),p.Post_ID)
	checkErr(err)

	return err
}

/* addPost saves the values of a new category to the database and is a method to Post.
 * Called by CreatePost
 * Connect to the DB and prepares query.
 * Execute query with the inserted values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (p *Post) save() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO posts (title,description,content,category_id) VALUES(?,?,?,?)")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description,[]byte(p.Content),p.Category_ID)
	fmt.Println(res)
	checkErr(err)
	return err
}
// End Post methods

/* EditPost takes updated form values from the http.request to populate a Post and call the savePost method.
 * The request delivers the FormValues if asked.
 * Convert post_id to an INT. The post ID is pulled from the from as a string.
 * FormValues are appointed to to the memory address of the Post struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call savePost, a method of the Post Struct, to update the DB
*/
func (post *Post) Patch(r *http.Request) (*Data,bool) {
	data := new(Data)
	updated := false

	if r.Method == "POST" {
		category_id := r.FormValue("Category_ID")
		category := r.FormValue("Category")
		/* 	To add a new category from a add post form we need to create a new
			 category, and then get the new ID of that category to insert it into the Post struct.
			 Also see addCategoryFromForm
		 */
		if (category != "") {
			category_id = addCategoryFromForm(category, category_id);
		} else {
			fmt.Println("empty string")
		}
		//// convert string values to INT before inserting into Struct and DB
		categoryINT, _ := strconv.Atoi(category_id)
		err := r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(post, r.PostForm)

		checkErr(err)

		post.Category_ID = categoryINT;
		if post.Title == "" || post.Content == "" {
			fmt.Println("I got here! ")
			data.Posts = append(data.Posts , post)
			data.Message = "Please fill in all the required fields"
			fmt.Println(data.Message)
		} else {
			err = post.update()
			checkErr(err)
			data.Posts = append(data.Posts , post)
			updated = true
		}
	} else {
		data.Posts = append(data.Posts , post)
	}
	return data,updated
}

/* NewPost takes updated form values from the http.request to populate a Post and call the addPost method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the Post struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call addPost, a method of the Post Struct, to insert new post in the DB.
*/
func Create(r *http.Request) (*Data,bool){

	data := new(Data)
	post := new(Post)

	created := false

	if r.Method == "POST" {
		category_id := r.FormValue("Category_ID")
		category := r.FormValue("Category")
		/* 	To add a new category from a add post form we need to create a new
			 category, and then get the new ID of that category to insert it into the Post struct.
			 Also see addCategoryFromForm
		 */
		if (category != "") {
			category_id = addCategoryFromForm(category, category_id);
		} else {
			fmt.Println("empty string")
		}
		//// convert string values to INT before inserting into Struct and DB
		categoryINT, _ := strconv.Atoi(category_id)
		err := r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(post, r.PostForm)
		checkErr(err)

		post.Category_ID = categoryINT;
		fmt.Printf("%v",post.Title)
		if post.Title == "" || post.Content == "" {
			data.Posts = append(data.Posts , post)
			data.Message = "Please fill in all the required fields"
		} else {
			err = post.save()
			checkErr(err)
			data.Posts = append(data.Posts , post)
			created = true
		}
	} else {
		data.Posts = append(data.Posts , post)
	}
	return data,created
}

/* addCategoryFromForm uses cat.AddCategory to add a new category.
 * Because a new category is created it is needed to get the ID from the database
  after creation so it can be added to returned and added inside a new or existing post.
*/
func addCategoryFromForm (category string, category_id string) string {
	//fmt.Print(category)
	c := &cat.Category{Title: category}
	fmt.Println(c)
	err := c.Store()
	checkErr(err)

	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)

	row := db.QueryRow("SELECT category_id FROM categories WHERE title = ? LIMIT 1",category)
	err = row.Scan(&category_id)
	checkErr(err)
	//fmt.Println("category_id: ",category_id)

	return category_id
}
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
