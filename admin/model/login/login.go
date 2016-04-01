package login

import (
"net/http"
"html/template"
"golang.org/x/crypto/bcrypt"
"database/sql"
_"github.com/go-sql-driver/mysql"
"github.com/jschalkwijk/GolangBlog/admin/config"
"github.com/gorilla/sessions"
"github.com/nu7hatch/gouuid"
)

var view = "GolangBlog/admin/view"
var templates = "GolangBlog/admin/templates"

type User struct {
	Username interface {}
	FirstName interface {}
	LastName interface {}
	Rights interface {}
	Logged interface{}

}
func RenderTemplate(w http.ResponseWriter,name string) {
	t, err := template.ParseFiles(templates+"/"+"header.html",view + "/" + name + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,name,nil)
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ExampleNewV4() string {
	newID, err := uuid.NewV4()
	checkErr(err)
	id := newID.String()
	return id
}

var store = sessions.NewCookieStore([]byte(ExampleNewV4()))

func SetSession(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	// Set some session values.
	session.Values["username"] = username
	session.Values["first-name"] = firstName
	session.Values["last-name"] = lastName
	session.Values["rights"] = rights
	session.Values["logged-in"] = true

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
}

func GetSession(r *http.Request) *User {
	session, err := store.New(r, "session")
	checkErr(err)
	user := &User{session.Values["username"],session.Values["first-name"],session.Values["last-name"],session.Values["rights"], session.Values["logged-in"]}
	return user
}

var username string
var hashedPassword []byte
var firstName string
var lastName string
var rights string

func Login(w http.ResponseWriter, r *http.Request) error {
	db, err := sql.Open("mysql",config.DB)
	checkErr(err)
	username = r.FormValue("username")
	password := []byte(r.FormValue("password"))
	row := db.QueryRow("SELECT username,password,first_name,last_name,rights FROM users WHERE username = ? LIMIT 1", username)
	row.Scan(&username,&hashedPassword,&firstName,&lastName,&rights)
	correct := compareHash(hashedPassword,password)
	return correct
}

func Logout(w http.ResponseWriter,r *http.Request){
	session, err := store.New(r, "session")
	checkErr(err)
	session.Options.MaxAge = -1
	session.Save(r, w)
}
func compareHash(hashedPassword []byte,password []byte) error{
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
