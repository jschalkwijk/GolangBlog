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
)

// here we define the absolute path to the view folder it takes the go root until the github folder.
var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

/* Post struct will hold data about a post and can be added to the Data struct */
type Post struct {
	Post_ID int
	Title string
	Description string
	Content template.HTML
	Keywords string
	Approved int
	Author string
	Date string
	Category_ID int
	Category string
	Trashed int
}

/*
 * Declaring vars corresponding to the struct. When scanning data from the database, the
   data will be stored on the memory address of these vars.
*/
var post_id int
var title string
var description string
var content string
var keywords string
var approved int
var author string
var date string
var category_id int
var category string
var trashed int

/* Stores a single post, or multiple posts in a Slice which can be iterated over in the template */
type Data struct {
	Posts      []Post
	Categories []cat.Category
	Deleted	bool
	Dashboard bool
}


/* RenderTemplate parses templates in cache before executing them. takes Response, content template name, and a Data struct
 * 	The function template.ParseFiles will read the contents of multiple "name".html files into cache.
 *	The method t.Execute executes the template, the string must correspond to the name giving to the template
 *	when defining them.
 *	After executing all the subtemplates, t.Execute will write the generated HTML to the http.ResponseWriter.
*/
func RenderTemplate(w http.ResponseWriter,name string, p *Data) {
	t, err := template.ParseFiles(templates+"/"+"header.html",templates+"/"+"nav.html",view + "/" + name + ".html",templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,p)
	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
func GetPosts(trashed int) *Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	// Selects all rows from posts, and links the category_id row to the matching title.
	rows, err := db.Query("SELECT posts.*, categories.title AS category FROM categories JOIN posts ON categories.categorie_id = posts.category_id WHERE posts.trashed = ? ORDER BY posts.post_id DESC",trashed)
	checkErr(err)

	data:= new(Data)

	for rows.Next() {
		err = rows.Scan(&post_id, &title, &description, &content,&keywords,&approved,
			&author,&date,&category_id,&trashed,&category)
		checkErr(err)
		// convert string to HTML markdown
		body := template.HTML(content)
		post := Post{post_id,title,description,body,keywords,approved,author,date,category_id,category,trashed}
		data.Posts = append(data.Posts , post)
		data.Dashboard = false
	}

	if(trashed == 1) {
		data.Deleted = true
	} else {
		data.Deleted = false
	}

	return data
}

/* -- Get a single Post -- */
/* GetSinglePost gets a post from the DB and returns a pointer to the Struct. It takes a id and post_title.
 * 	Connects to the database and gets all post rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new Post Struct and insert values.
 *	Append the Post struct to the Data.Posts Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func GetSinglePost(id string,post_title string, getCat bool) *Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection established")
	defer db.Close()
	defer fmt.Println("Connection Closed")

	rows := db.QueryRow("SELECT posts.*, categories.title AS category FROM categories JOIN posts ON categories.categorie_id = posts.category_id WHERE post_id=? LIMIT  1", id)

	collection := new(Data)

	err = rows.Scan(&post_id, &title, &description, &content,&keywords,&approved,&author,&date,&category_id,&trashed,&category)
	checkErr(err)

	body := template.HTML(content)
	post := Post{post_id,title,description,body,keywords,approved,author,date,category_id,category,trashed}

	collection.Posts = append(collection.Posts , post)
	 /* When we need to edit or create a post, we need to get the categories in order to select them inside the html page.
	  * since we already have a function inside the categories model, we will call that.
	  * This returns a pointer to the Data struct of model/categories. We  set our
	    categories.Data struct, Categories, to the slice of posts.Data.Categories.
	  * They are accessible inside the template now.
	 */
	if(getCat) {
		listCat := cat.GetCategories(0)
		collection.Categories = listCat.Categories
	}

	fmt.Println(collection.Categories)
	return collection
}

/* -- Post Methods -- */

/* savePost updates the values of an existing post to the database and is a method to Post
 * Called by EditPost
 * Connect to the DB and prepares query.
 * Execute query with the inserted struct values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (p *Post) savePost() error {
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
	res, err := stmt.Exec(p.Title,p.Description,p.Category_ID,[]byte(p.Content),p.Post_ID)
	checkErr(err)
	//affect, err := res.RowsAffected()
	//checkErr(err)
	//fmt.Println(affect)
	fmt.Println(res)
	return err
}

/* addPost saves the values of a new category to the database and is a method to Post.
 * Called by NewPost
 * Connect to the DB and prepares query.
 * Execute query with the inserted values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (p *Post) addPost() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO posts (title,description,content,category_id) VALUES(?,?,?,?)")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description,[]byte(p.Content),p.Category_ID)
	affect, err := res.RowsAffected()
	fmt.Println(affect)
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
func EditPost(w http.ResponseWriter, r *http.Request,id string,title string) {
	title = r.FormValue("title")
	description := r.FormValue("description")
	category_id := r.FormValue("selected-category")
	content := r.FormValue("content")

	category := r.FormValue("category")

	/* 	To add a new category from a edit post form we need to create a new
	 	category, and then get the new ID of that category to insert it into the Post struct.
	 	Also see addCategoryFromForm
	 */
	if (category != "") {
		category_id = addCategoryFromForm(category,category_id);
	} else {
		fmt.Println("empty string")
	}
	// convert string to INT before inserting into Struct and DB
	idINT,error := strconv.Atoi(id)
	checkErr(error)
	categoryINT,error := strconv.Atoi(category_id)
	checkErr(error)
	// Convert string to HTML
	body := template.HTML(content)
	p := &Post{Post_ID: idINT, Title: title,Description: description,Category_ID: categoryINT, Content: body}
	fmt.Println(p)
	err := p.savePost()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/posts/"+id+"/"+title, http.StatusFound)
}

/* NewPost takes updated form values from the http.request to populate a Post and call the addPost method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the Post struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call addPost, a method of the Post Struct, to insert new post in the DB.
*/
func NewPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	category_id := r.FormValue("selected-category")
	content := r.FormValue("content")
	category := r.FormValue("category")
	/* 	To add a new category from a add post form we need to create a new
	 	category, and then get the new ID of that category to insert it into the Post struct.
	 	Also see addCategoryFromForm
	 */
	if (category != "") {
		category_id = addCategoryFromForm(category,category_id);
	} else {
		fmt.Println("empty string")
	}
	// Convert string to HTML
	body := template.HTML(content)
	// convert string values to INT before inserting into Struct and DB
	categoryINT,error := strconv.Atoi(category_id)
	checkErr(error)
	p := &Post{Title: title ,Description: description, Content: body,Category_ID: categoryINT}
	fmt.Println(p)
	err := p.addPost()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/posts", http.StatusFound)
}

/* addCategoryFromForm uses cat.AddCategory to add a new category.
 * Because a new categpry is created it is needed to get the ID from the database
  after creation so it can be added to returned and added inside a new or existing post.
*/
func addCategoryFromForm (category string, category_id string) string {
	//fmt.Print(category)
	c := &cat.Category{Title: category}
	fmt.Println(c)
	err := c.AddCategory()
	checkErr(err)

	db, err := sql.Open("mysql", config.DB)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")
	checkErr(err)

	row := db.QueryRow("SELECT categorie_id FROM categories WHERE title = ? LIMIT 1",category)
	err = row.Scan(&category_id)
	checkErr(err)
	//fmt.Println("category_id: ",category_id)

	return category_id
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
