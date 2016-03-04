package data

import (

)

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


type Data struct {
	Posts []Post
}