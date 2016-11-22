package config

// Declare database variables
var user string = "root"
var password string= ""
var host string = "localhost"
var name string = "golangblog"
var port string = "3306"

//db_user,db_pass,@tcp(db_host:port)/db_name?charset=utf8
var DB string = user+":"+password+"@tcp("+host+":"+port+")/"+name+"?charset=utf8"

/* Declaring folder paths */
var View = "GolangBlog/admin/view"
var Templates = "GolangBlog/admin/templates"