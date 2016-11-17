package pages

import(
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"html/template"
	"net/http"
	"github.com/gorilla/schema"
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
	res,err := stmt.Exec(p.Title,p.Description,p.Content)
	if(res.RowsAffected() == 1){
		println("Page created")
	}
	return err
}

var page_id int
var title string
var description string
var content string
var keywords string
var approved bool
var author string
var date string
var parent_id int
var trashed bool

type Data struct {
	Pages	[]Page
}

func All(trashed bool) *Data {
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM pages WHERE trashed = ?",trashed)

	var page Page
	data := new(Data)

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

func New(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	checkErr(err)

	schema.NewDecoder()
	p := new(Page)
	decoder := schema.NewDecoder()
	err = decoder.Decode(p, r.PostForm)
	p.save();

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}