package categories
import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/model/categories"
	cat "github.com/jschalkwijk/GolangBlog/admin/model/categories"
	//"github.com/jschalkwijk/GolangBlog/controller"
	"github.com/gorilla/mux"
)


func Index(w http.ResponseWriter, r *http.Request) {
	p := cat.All(0)
	categories.View(w,"categories", p)
}

func One(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	p := cat.One(id)
	categories.View(w,"categories", p)
}