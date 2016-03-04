package db

//db_user,db_pass,@tcp(db_host:port)/db_name?charset=utf8
var user string = "root"
var password string= "root"
var host string = "localhost"
var name string = "nerdcms_db"

var DB string = user+":"+password+"@tcp("+host+":8889)/"+name+"?charset=utf8"
