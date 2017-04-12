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
	"log"
)

type User struct {
	Username string
	FirstName string
	LastName string
	Rights string
	Logged bool
}

var username string
var hashedPassword []byte
var firstName string
var lastName string
var rights string

func View(w http.ResponseWriter,name string) {
	t, err := template.ParseFiles(config.Templates+"/"+"header.html",config.View + "/" + name + ".html")
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
	session, err := store.Get(r,"session")
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
	session, _ := store.New(r,"session")
	if (session.Values["username"] != nil) {
		// Kan ik niet beter dit een keer doen en dan een referene naar de user struct geven of moet ik dit elke keer aanroepen?
		user := &User{session.Values["username"].(string), session.Values["first-name"].(string), session.Values["last-name"].(string), session.Values["rights"].(string), session.Values["logged-in"].(bool)}
		return user
	} else {
		user := new(User)
		return user
	}
}


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
		log.Println(err)
	}
}
