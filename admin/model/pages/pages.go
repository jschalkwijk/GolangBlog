package pages

import(
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"html/template"
	"net/http"
	"github.com/gorilla/schema"
	"fmt"
)

type Page struct {
	Page_ID int `schema:"-"`
	Title string `schema:"title"`
	Description string `schema:"description"`
	Content template.HTML `schema:"content"`
	Keywords string `schema:"-"`
	Approved int `schema:"-"`
	Author string `schema:"-"`
	Date string `schema:"-"`
	Parent_ID int `schema:"category_id"`
	Trashed int `schema:"-"`
}

func (p *Page) save() error {
	db ,err := sql.Open("mysql", config.DB)
	checkErr(err)

	stmt,err := db.Prepare("INSERT INTO pages(title,description,content) VALUES(?,?,?)")
	checkErr(err)
	_,err = stmt.Exec(p.Title,p.Description,[]byte(p.Content))
	checkErr(err)

	return err
}

func (p *Page) update() error {
	db ,err := sql.Open("mysql", config.DB)
	checkErr(err)

	stmt,err := db.Prepare("UPDATE pages SET title = ?, description = ?,content = ? WHERE page_id = ?")
	checkErr(err)
	_,err = stmt.Exec(p.Title,p.Description,[]byte(p.Content),p.Page_ID)
	checkErr(err)

	return err
}

var page_id int
var title string
var description string
var content string
var keywords string
var approved int
var author string
var date string
var parent_id int
var trashed int

type Data struct {
	Pages	[]Page
	Message string
}

func All(trashed int) *Data {
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM pages WHERE trashed = ? ORDER BY page_id DESC",trashed)

	var page Page
	data := new(Data)
	/*
		data := new(Data)
		page := new(Page)
		err = rows.Scan(&page.Page_ID,&title,&description,&content,&keywords,&approved,&author,&date,&parent_id,&trashed)
	*/
	for rows.Next() {
		err = rows.Scan(&page_id,&title,&description,&content,&keywords,&approved,&author,&date,&parent_id,&trashed)
		checkErr(err)
		body := template.HTML(content)
		page = Page{
			page_id,
			title,
			description,
			body,
			keywords,
			approved,
			author,
			date,
			parent_id,
			trashed,
		}

		data.Pages = append(data.Pages , page)
	}
	return data
}

func Single(id string) *Data {
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()
	rows := db.QueryRow("SELECT * FROM pages WHERE page_id = ?",id)

	var page Page
	data := new(Data)

	err = rows.Scan(&page_id,&title,&description,&content,&keywords,&approved,&author,&date,&parent_id,&trashed)
	checkErr(err)
	body := template.HTML(content)
	page = Page{
		page_id,
		title,
		description,
		body,
		keywords,
		approved,
		author,
		date,
		parent_id,
		trashed,
	}
	data.Pages = append(data.Pages , page)
	return data
}

func Create(r *http.Request) (*Data,bool){
	data := new(Data)
	page := new(Page)

	created := false

	if r.Method == "POST" {
		err := r.ParseForm()
		checkErr(err)
		// new schema to pass form values to the Page struct
		schema.NewDecoder()
		decoder := schema.NewDecoder()
		err = decoder.Decode(page, r.PostForm)
		checkErr(err)
		fmt.Println("schema: ",page)
		if page.Title == "" || page.Description == "" || page.Content == "" {
			data.Pages = append(data.Pages , *page)
			data.Message = "Please fill in all the required fields"
		} else {
			err = page.save()
			checkErr(err)
			data.Pages = append(data.Pages , *page)
			created = true
		}
	} else {
		page := new(Page)
		data.Pages = append(data.Pages , *page)
	}

	return data,created
}

func Patch(page Page,r *http.Request) (*Data,bool){
	data := new(Data)
	data.Pages = append(data.Pages , page)
	updated := false
	// Ik moet van Page een pointer maken naar Page anders werkt het niet zoals ik wil.
	if r.Method == "POST" {
		err := r.ParseForm()
		checkErr(err)
		// new schema to pass form values to the Page struct
		schema.NewDecoder()
		decoder := schema.NewDecoder()
		err = decoder.Decode(data.Pages[0], r.PostForm)
		checkErr(err)
		fmt.Println("schema: ",page)
		if page.Title == "" || page.Description == "" || page.Content == "" {
			data.Pages = append(data.Pages , page)
			data.Message = "Please fill in all the required fields"
		} else {
			err = page.update()
			checkErr(err)
			data.Pages = append(data.Pages , page)
			updated = true
		}
	} else {
		data.Pages = append(data.Pages , page)
	}

	return data,updated
}
//func (p *Page) Patch(r *http.Request) (*Data,bool){
//	data := new(Data)
//	data.Pages = append(data.Pages , p)
//	updated := false
//	// Ik moet van Page een pointer maken naar Page anders werkt het niet zoals ik wil.
//	if r.Method == "POST" {
//		err := r.ParseForm()
//		checkErr(err)
//		// new schema to pass form values to the Page struct
//		schema.NewDecoder()
//		decoder := schema.NewDecoder()
//		err = decoder.Decode(data.Pages[0], r.PostForm)
//		checkErr(err)
//		fmt.Println("schema: ",p)
//		if p.Title == "" || p.Description == "" || p.Content == "" {
//			data.Pages = append(data.Pages , p)
//			data.Message = "Please fill in all the required fields"
//		} else {
//			err = p.update()
//			checkErr(err)
//			data.Pages = append(data.Pages , p)
//			updated = true
//		}
//	} else {
//		data.Pages = append(data.Pages , p)
//	}
//
//	return data,updated
//}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}