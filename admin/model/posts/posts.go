package posts

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"github.com/jschalkwijk/GolangBlog/admin/model/db"
)

// here we define the absolute path to the view folder it takes the go root until the github folder.
var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

// Post struct to create blog which will be added to the collection struct
type Post struct {
	Post_ID int
	Title string
	Description string
	Content string
	Keywords string
	Approved int
	Author string
	Date string
	Category_ID int
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
var trashed int

// Stores a single post, or multiple blog which we can then iterate over in the template
type Data struct {
	Posts []Post
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
	db, err := sql.Open("mysql", db.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	rows, err := db.Query("SELECT * FROM posts ORDER BY post_id DESC")
	checkErr(err)

	//collection := Posts
	collection := new(Data)
	for rows.Next() {
		err = rows.Scan(&post_id, &title, &description, &content,&keywords,&approved,
			&author,&date,&category_id,&trashed)
		checkErr(err)

		post := Post{post_id,title,description,content,keywords,approved,author,date,category_id,trashed}

		collection.Posts = append(collection.Posts , post)
	}

	return collection
}

//Get a single Post
func GetSinglePost(id string,post_title string) *Data {
//	vars := mux.Vars(r)
//	id:= vars["id"]
//	post_title := vars["title"]
	db, err := sql.Open("mysql", db.DB)
	checkErr(err)
	fmt.Println("Connection established")
	defer db.Close()
	defer fmt.Println("Connection Closed")
	fmt.Println("SELECT * FROM posts WHERE post_id="+id+" AND title='"+post_title+"' LIMIT  1")
	rows := db.QueryRow("SELECT * FROM posts WHERE post_id=? AND title=? LIMIT 1", id,post_title)

	collection := new(Data)

	err = rows.Scan(&post_id, &title, &description, &content,&keywords,&approved,
		&author,&date,&category_id,&trashed)
	checkErr(err)

	post := Post{post_id,title,description,content,keywords,approved,author,date,category_id,trashed}

	collection.Posts = append(collection.Posts , post)

	//fmt.Println(collection.Posts)
	return collection
}

// Post Methods
func (p *Post) savePost() error {
	db, err := sql.Open("mysql", db.DB)
	defer db.Close()
	checkErr(err)
	stmt, err := db.Prepare("UPDATE posts SET title=?, description=?, content=? WHERE post_id=?")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description,p.Content,p.Post_ID)
	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	fmt.Println(res)
	return err
}

func (p *Post) addPost() error {
	db, err := sql.Open("mysql", db.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO posts (title,description,content) VALUES(?,?,?) ")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description,p.Content)
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
	//category_id := r.FormValue("category_id")
	content := r.FormValue("content")
	new_id,error := strconv.Atoi(id)
	checkErr(error)
	p := &Post{Post_ID: new_id, Title: title,Description: description, Content: content}
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
	//category_id := r.FormValue("category_id")
	content := r.FormValue("content")

	p := &Post{Title: title ,Description: description, Content: content}
	fmt.Println(p)
	err := p.addPost()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/posts/", http.StatusFound)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
