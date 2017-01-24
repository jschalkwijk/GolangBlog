package pages

import(
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"html/template"
	"net/http"
	"github.com/gorilla/schema"
	"fmt"
	"log"
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

func (p *Page) Patch(r *http.Request) (*Data,bool){
	data := new(Data)
	updated := false

	if r.Method == "POST" {
		err := r.ParseForm()
		checkErr(err)
		println(r.FormValue("title"))
		// new schema to pass form values to the Page struct
		schema.NewDecoder()
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(p, r.PostForm)
		checkErr(err)
		fmt.Println("schema: ",p)
		if p.Title == "" || p.Description == "" || p.Content == "" {
			data.Pages = append(data.Pages , p)
			data.Message = "Please fill in all the required fields"
		} else {
			err = p.update()
			checkErr(err)
			data.Pages = append(data.Pages , p)
			updated = true
		}
	} else {
		data.Pages = append(data.Pages , p)
	}

	return data,updated
}

var content string

type Data struct {
	Pages	[]*Page
	Message string
}

func All(trashed int) *Data {
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM pages WHERE trashed = ? ORDER BY page_id DESC",trashed)

	data := new(Data)

	for rows.Next() {
		page := new(Page)
		err = rows.Scan(&page.Page_ID,&page.Title,&page.Description,&content,&page.Keywords,&page.Approved,&page.Author,&page.Date,&page.Parent_ID,&page.Trashed)
		checkErr(err)
		page.Content = template.HTML(content)

		data.Pages = append(data.Pages , page)
	}
	return data
}

func One(id string) *Data {
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()
	rows := db.QueryRow("SELECT * FROM pages WHERE page_id = ?",id)

	page := new(Page)
	data := new(Data)

	err = rows.Scan(&page.Page_ID,&page.Title,&page.Description,&content,&page.Keywords,&page.Approved,&page.Author,&page.Date,&page.Parent_ID,&page.Trashed)
	checkErr(err)
	page.Content = template.HTML(content)
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
			data.Pages = append(data.Pages , page)
			data.Message = "Please fill in all the required fields"
		} else {
			err = page.save()
			checkErr(err)
			data.Pages = append(data.Pages , page)
			created = true
		}
	} else {
		page := new(Page)
		data.Pages = append(data.Pages , page)
	}

	return data,created
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}