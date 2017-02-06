package users

import (
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
)

type User struct {
	User_ID int `schema:"-"`
	UserName string
	Password string
	First_Name string
	Last_Name string
	DOB string
	Email string
	Function string
	Rights string
	Approved int
	Trashed int
}

type Data struct {
	Users []*User
	Deleted	bool
	Message string
}

/* -- Get all Users --
 * 	Connects to the database and gets all posts rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	For every row get the values, and set the values to the memory address of the named variable.
 		- Instantiate a new User Struct and insert values.
 		- Append the User struct to the Data.Users Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func All(trashed int) *Data {
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	fmt.Println("Connection with database Established")
	defer db.Close()
	defer fmt.Println("Connection with database Closed")

	// Selects all rows from posts, and links the category_id row to the matching title.
	rows, err := db.Queryx("SELECT * FROM users WHERE trashed = ? ORDER BY user_id DESC",trashed)
	checkErr(err)

	data := new(Data)

	for rows.Next() {
		user := new(User)
		err = rows.StructScan(
			&user,
		)
		checkErr(err)

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
/* GetOneUser gets a user from the DB and returns a pointer to the Struct. It takes a id and username.
 * 	Connects to the database and gets all post rows.
 * 	Instantiate a new Data struct assigned to var collection
 * 	Get a single row from the DB and get the values and set the values to the memory address of the named variable.
 *	Instantiate a new User Struct and insert values.
 *	Append the Post struct to the Data.Users Slice.
 *	Returns the Data Struct after the loop is completed. This Struct can be used
  	inside a template.
 */
func One(id string) *Data {
	db, err := sqlx.Connect("mysql", config.DB)
	checkErr(err)
	defer db.Close()

	rows := db.QueryRowx("SELECT * FROM users WHERE user_id = ?", id)

	data := new(Data)

	user := new(User)

	err = rows.StructScan(
		&user,
	)
	checkErr(err)

	data.Users = append(data.Users , user)

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
func (u *User) update() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	checkErr(err)

	fmt.Println("reference to User struct: ", u)

	stmt, err := db.Prepare("UPDATE users SET username=?,password=?, first_name=?, last_name=?, dob=?, email=?, function=?, rights=? WHERE user_id=?")
	fmt.Println(stmt)
	checkErr(err)

	res, err := stmt.Exec(u.UserName,u.Password,u.First_Name,u.Last_Name,u.DOB,u.Email,u.Function,u.Rights,u.User_ID)
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
func (u *User) store() error {
	db, err := sql.Open("mysql", config.DB)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO users (username,password, first_name, last_name, dob, email, function, rights) VALUES(?,?,?,?,?,?,?,?)")
	fmt.Println(stmt)
	checkErr(err)
	res, err := stmt.Exec(u.UserName,u.Password,u.First_Name,u.Last_Name,u.DOB,u.Email,u.Function,u.Rights)
	affect, err := res.RowsAffected()
	fmt.Println(affect)
	fmt.Println(res)
	checkErr(err)
	return err
}
// End User methods

func (user *User) Patch(r *http.Request) (*Data,bool) {
	data := new(Data)
	updated := false

	if r.Method == "POST" {
		err := r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(user, r.PostForm)
		checkErr(err)

		passwordCheck := r.FormValue("Password-again");
		if (user.Password == passwordCheck ){
			user.Password = hashPassword([]byte(passwordCheck))
		} else {
			fmt.Println("Password not set")
		}
		fmt.Println(user)

		if user.UserName == "" || user.Email == "" {
			data.Users = append(data.Users , user)
			data.Message = "Please fill in all the required fields"
		} else {
			err = user.update()
			checkErr(err)
			data.Users = append(data.Users , user)
			updated = true
		}
	} else {
		data.Users = append(data.Users , user)
	}

	return data,updated
}

/* NewUser takes updated form values from the http.request to populate a User strut and call the addUser method.
 * The request delivers the FormValues if asked.
 * FormValues are appointed to to the memory address of the User struct. There is only one to edit so no need to
   instantiate a separate one.
 * Call addUser, a method of the User Struct, to insert new user in the DB.
*/
func Create(r *http.Request) (*Data,bool) {
	user := new(User)
	data := new(Data)

	created := false

	if r.Method == "POST" {
		err := r.ParseForm()
		decoder := schema.NewDecoder()
		decoder.ZeroEmpty(true)
		err = decoder.Decode(user, r.PostForm)
		checkErr(err)

		if user.UserName == "" || user.Email == "" || user.Password == "" || r.FormValue("Password-again") == "" {
			data.Users = append(data.Users , user)
			data.Message = "Please fill in all the required fields"
		} else {
			passwordCheck := r.FormValue("Password-again");
			if (user.Password == passwordCheck ){
				user.Password = hashPassword([]byte(passwordCheck))
			} else {
				fmt.Println("Password not set")
			}
			fmt.Println(user)
			err = user.store()
			data.Users = append(data.Users , user)
			created = true
		}
	} else {
		data.Users = append(data.Users , user)
	}
	return data,created
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

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}