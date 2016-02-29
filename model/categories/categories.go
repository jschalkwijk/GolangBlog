package categories

import (
_"github.com/go-sql-driver/mysql"
"database/sql"
"fmt"
"html/template"
"net/http"
"path/filepath"
"strconv"
)

// here we define the absolute path to the view folder it takes the go root until the github folder.
var view, _ = filepath.Abs("../jschalkwijk/GolangBlog/view")
var templates, _ = filepath.Abs("../jschalkwijk/GolangBlog/templates")

// categorie struct to create categories which will be added to the collection struct
type Categorie struct {
	Category_ID int
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

var category_id int
var title string
var description string
var content string
var keywords string
var approved int
var author string
var cat_type string
var date string
var parent_id int
var trashed int

// Stores a single categorie, or multiple categories which we can then iterate over in the template
type Data struct {
	Categories []Categorie
}

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

// Get all categories
func GetCategories() *Data {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/nerdcms_db?charset=utf8")
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	rows, err := db.Query("SELECT * FROM categories ORDER BY categorie_id DESC")
	checkErr(err)

	collection := new(Data)
	for rows.Next() {
		err = rows.Scan(&category_id, &title, &description, &content,&keywords,&approved,
		&author,&cat_type,&date,&parent_id,&trashed)
		checkErr(err)

		category := Categorie{category_id,title,description,content,keywords,approved,author,cat_type,date,parent_id,trashed}

		collection.Categories = append(collection.Categories , category)
	}

	return collection
}

//Get a single categorie
func GetSingleCategory(id string,category_title string) *Data {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/nerdcms_db?charset=utf8")
	checkErr(err)
	fmt.Println("Connection established")
	defer db.Close()
	defer fmt.Println("Connection Closed")
	fmt.Println("SELECT * FROM categories WHERE categorie_id="+id+" AND title='"+category_title+"' LIMIT  1")
	rows := db.QueryRow("SELECT * FROM categories WHERE categorie_id=? AND title=? LIMIT 1", id,category_title)

	collection := new(Data)

	err = rows.Scan(&category_id, &title, &description, &content,&keywords,&approved,
	&author,&cat_type,&date,&parent_id,&trashed)
	checkErr(err)

	category := Categorie{category_id,title,description,content,keywords,approved,author,cat_type,date,parent_id,trashed}

	collection.Categories = append(collection.Categories , category)

	//fmt.Println(collection.categories)
	return collection
}

// categorie Methods
func (p *Categorie) saveCategory() error {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/nerdcms_db?charset=utf8")
	defer db.Close()
	checkErr(err)
	stmt, err := db.Prepare("UPDATE categories SET title=?, description=? WHERE categorie_id=?")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description,p.Category_ID)
	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	fmt.Println(res)
	return err
}

func (p *Categorie) addCategory() error {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/nerdcms_db?charset=utf8")
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO categories (title,description) VALUES(?,?) ")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(p.Title,p.Description)
	affect, err := res.RowsAffected()
	fmt.Println(affect)
	fmt.Println(res)
	checkErr(err)
	return err
}
// End category methods


func EditCategory(w http.ResponseWriter, r *http.Request,id string,title string) {
	title = r.FormValue("title")
	description := r.FormValue("description")
	id_string := r.FormValue("category_id")
	category_id,error := strconv.Atoi(id_string)
	checkErr(error)
	p := &Categorie{Category_ID: category_id, Title: title,Description: description}
	fmt.Println(p)
	err := p.saveCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories/"+id+"/"+title, http.StatusFound)
}

func NewCategory(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")

	p := &Categorie{Title: title ,Description: description}
	fmt.Println(p)
	err := p.addCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusFound)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

