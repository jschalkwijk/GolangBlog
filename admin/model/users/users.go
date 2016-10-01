package users

import (
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"fmt"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID int
	UserName string
	Password string
	FirstName string
	LastName string
	DOB string
	Email string
	Function string
	Rights string
	Approved int
	Trashed int
}

type Data struct {
	Users []User
	Deleted	bool
}

/*
 * Declaring vars corresponding to the struct. When scanning data from the database, the
   data will be stored on the memory address of these vars.
*/
var userID int
var username string
var password string
var firstName string
var lastName string
var email string
var dob string
var function string
var rights string
var trashed int
var approved int

/* -- Get all Users --
 * 	Connects to the database and gets all posts rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	For every row get the values, and set the values to the memory address of the named variable.
 		- Instantiate a new User Struct and insert values.
 		- Append the User struct to the Data.Users Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func (data Data)Get(trashed int) Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	// Selects all rows from posts, and links the category_id row to the matching title.
	rows, err := db.Query("SELECT * FROM users WHERE trashed = ? ORDER BY user_id DESC",trashed)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&userID, &username, &password, &firstName,&lastName,
			&dob,&email,&function,&rights,&trashed,&approved)
		checkErr(err)
		// convert string to HTML markdown
		user := User{userID, username, password, firstName,lastName,dob,email,
			function,rights,trashed,approved}
		data.Users = append(data.Users, user)
	}

	if(trashed == 1) {
		data.Deleted = true
	} else {
		data.Deleted = false
	}

	return data
}

/* -- Get a single User -- */
/* GetSingleUser gets a user from the DB and returns a pointer to the Struct. It takes a id and username.
 * 	Connects to the database and gets all post rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new User Struct and insert values.
 *	Append the Post struct to the Data.Users Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func GetSingleUser(id string,post_title string) *Data {
	db, err := sql.Open("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection established")
	defer db.Close()
	defer fmt.Println("Connection Closed")

	rows := db.QueryRow("SELECT * FROM users WHERE user_id = ?", id)

	data := new(Data)

	err = rows.Scan(&userID, &username, &password, &firstName,&lastName,&email,
		&dob,&function,&rights,&trashed,&approved)
	checkErr(err)

	user := User{userID, username, password, firstName,lastName,email,
		dob,function,rights,trashed,approved}

	data.Users = append(data.Users , user)
	err = compareHash([]byte(user.Password))
	checkErr(err)
	return data
}

/* -- User Methods -- */

/* saveUser updates the values of an existing post to the database and is a method to User
 * Called by EditUser
 * Connect to the DB and prepares query.
 * Execute query with the inserted struct values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (u *User) saveUser() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)

	fmt.Println("reference to User struct: ", u)

	stmt, err := db.Prepare("UPDATE users SET username=?,password=?, first_name=?, last_name=?, dob=?, email=?, function=?, rights=? WHERE user_id=?")
	fmt.Println(stmt)
	checkErr(err)

	res, err := stmt.Exec(u.UserName,u.Password,u.FirstName,u.LastName,u.DOB,u.Email,u.Function,u.Rights,u.UserID)
	checkErr(err)
	//affect, err := res.RowsAffected()
	//checkErr(err)
	//fmt.Println(affect)
	fmt.Println(res)
	return err
}

/* addUser saves the values of a new category to the database and is a method to User.
 * Called by NewUser
 * Connect to the DB and prepares query.
 * Execute query with the inserted values and replaces the ? in the query string.
 * Checks how many rows are affected.
 * Returns an error if needed.
*/
func (u *User) addUser() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO users (username,password, first_name, last_name, dob, email, function, rights) VALUES(?,?,?,?,?,?,?,?)")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(u.UserName,u.Password,u.FirstName,u.LastName,u.DOB,u.Email,u.Function,u.Rights)
	affect, err := res.RowsAffected()
	fmt.Println(affect)
	fmt.Println(res)
	checkErr(err)
	return err
}
// End User methods

func EditUser(w http.ResponseWriter, r *http.Request,id string,username string) {
	username = r.FormValue("username")
	newPassword := r.FormValue("new-password")
	checkPassword := r.FormValue("new-password-again")
	firstName := r.FormValue("first-name")
	lastName := r.FormValue("last-name")
	dob := r.FormValue("dob")
	email := r.FormValue("email")
	function := r.FormValue("function")
	rights := r.FormValue("rights")

	var password string

	if (newPassword == checkPassword){
		password = newPassword
	}
	// convert string to INT before inserting into Struct and DB
	idINT,error := strconv.Atoi(id)
	checkErr(error)
	u:= &User{UserID: idINT,UserName: username,Password: password, FirstName: firstName,LastName: lastName,DOB: dob,Email: email,Function: function,Rights: rights}
	fmt.Println(u)
	err := u.saveUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/users/"+id+"/"+username, http.StatusFound)
}

/* NewUser takes updated form values from the http.request to populate a User strut and call the addUser method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the User struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call addUser, a method of the User Struct, to insert new user in the DB.
*/
func NewUser(w http.ResponseWriter, r *http.Request) {
	username = r.FormValue("username")
	newPassword := r.FormValue("new-password")
	checkPassword := r.FormValue("new-password-again")
	firstName := r.FormValue("first-name")
	lastName := r.FormValue("last-name")
	dob := r.FormValue("dob")
	email := r.FormValue("email")
	function := r.FormValue("function")
	rights := r.FormValue("rights")

	var password string

	if (newPassword == checkPassword){
		password = hashPassword([]byte(newPassword))
	} else {
		fmt.Println("Password not set")
	}
	u := &User{UserName: username,Password: password, FirstName: firstName,LastName: lastName,DOB: dob,Email: email,Function: function,Rights: rights}
	fmt.Println(u)
	err := u.addUser()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func hashPassword(password []byte) (string){

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 15)
	if err != nil {
	panic(err)
	}
	fmt.Println(string(hashedPassword))
	return string(hashedPassword)
}

func compareHash(hashedPassword []byte) error{
		password := []byte("root")
		// Comparing the password with the hash
		err := bcrypt.CompareHashAndPassword(hashedPassword, password)
		return err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}