package posts

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"github.com/jschalkwijk/GolangBlog/admin/model/config"
	cat "github.com/jschalkwijk/GolangBlog/admin/model/categories"
)

// here we define the absolute path to the view folder it takes the go root until the github folder.
var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

// Post struct to create blog which will be added to the collection struct
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

// Stores a single post, or multiple blog which we can then iterate over in the template
type Data struct {
	Posts      []Post
	Categories []cat.Category
}
//var Posts []Post

/*
  The function template.ParseFiles will read the contents of "".html and return a *template.Template.
  The method t.Execute executes the template, writing the generated HTML to the http.ResponseWriter.
  The .Title and .Body dotted identifiers inside the template refer to p.Title and p.Body.
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

// Get all Posts
func GetPosts() *Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")


	rows, err := db.Query("SELECT posts.*, categories.title AS category FROM categories JOIN posts ON categories.categorie_id = posts.category_id order by posts.post_id DESC")
//	err = config.QueryRow("SELECT categories.title as cat FROM categories JOIN posts ON categories.categorie_id = posts.category_id").Scan(&cat)
	checkErr(err)
//	fmt.Println(cat)
	collection := new(Data)
	for rows.Next() {
		err = rows.Scan(&post_id, &title, &description, &content,&keywords,&approved,
			&author,&date,&category_id,&trashed,&category)
		checkErr(err)
		// convert string to HTML markdown
		body := template.HTML(content)
		post := Post{post_id,title,description,body,keywords,approved,author,date,category_id,category,trashed}
		collection.Posts = append(collection.Posts , post)
	}

	return collection
}

//Get a single Post
func GetSinglePost(id string,post_title string, getCat bool) *Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection established")
	defer db.Close()
	defer fmt.Println("Connection Closed")
	fmt.Println("SELECT * FROM posts WHERE post_id="+id+" AND title='"+post_title+"' LIMIT  1")
	rows := db.QueryRow("SELECT posts.*, categories.title AS category FROM categories JOIN posts ON categories.categorie_id = posts.category_id WHERE post_id=? LIMIT  1", id)

	collection := new(Data)

	err = rows.Scan(&post_id, &title, &description, &content,&keywords,&approved,&author,&date,&category_id,&trashed,&category)
	checkErr(err)

	body := template.HTML(content)
	post := Post{post_id,title,description,body,keywords,approved,author,date,category_id,category,trashed}

	collection.Posts = append(collection.Posts , post)

	// When we need to edit or create a post, we need to get the categories in order to select them inside the html page.
	// since we already have a function inside the categories model, we will call that.
	// This returns a pointer to the Data struct of model/categories. We  set our
	// Data struct, categories, to the slice of categories.
	if(getCat) {
		test := cat.GetCategories()
		collection.Categories = test.Categories
	}

	fmt.Println(collection.Categories)
	return collection
}

// Post Methods
func (p *Post) savePost() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)

	fmt.Println("reference to Post struct: ", p)

	stmt, err := db.Prepare("UPDATE posts SET title=?, description=?,category_id=?, content=? WHERE post_id=?")
	fmt.Println(stmt)
	checkErr(err)
	// to be able to save the new html to the database, we have to convert it to a slice of bytes, why is this working?, we can't save
	// a value of type template.HTML to the DB. I tried different things, change the .Content to string, byte, but then I have a problem displaying
	// the content in html format on the page.
	res, err := stmt.Exec(p.Title,p.Description,p.Category_ID,[]byte(p.Content),p.Post_ID)
	checkErr(err)
	//affect, err := res.RowsAffected()
	//checkErr(err)

	//fmt.Println(affect)
	fmt.Println(res)
	return err
}

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


func EditPost(w http.ResponseWriter, r *http.Request,id string,title string) {
	title = r.FormValue("title")
	description := r.FormValue("description")
	category_id := r.FormValue("selected-category")
	content := r.FormValue("content")

	category := r.FormValue("category")

	if (category != "") {
		category_id = addCategoryFromForm(category,category_id);
	} else {
		fmt.Println("empty string")
	}

	idINT,error := strconv.Atoi(id)
	checkErr(error)
	categoryINT,error := strconv.Atoi(category_id)
	checkErr(error)
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
func NewPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	category_id := r.FormValue("selected-category")
	content := r.FormValue("content")
	category := r.FormValue("category")

	if (category != "") {
		category_id = addCategoryFromForm(category,category_id);
	} else {
		fmt.Println("empty string")
	}

	body := template.HTML(content)
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
	fmt.Println("category_id: ",category_id)

	return category_id
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
