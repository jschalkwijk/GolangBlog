package categories
import (
	"net/http"
	"github.com/jschalkwijk/GolangBlog/model/categories"
	cat "github.com/jschalkwijk/GolangBlog/admin/model/categories"
	//"github.com/jschalkwijk/GolangBlog/controller"
	"github.com/gorilla/mux"
)


func Index(w http.ResponseWriter, r *http.Request) {
	p := cat.GetCategories(0)
	categories.RenderTemplate(w,"categories", p)
}

func Single(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	post_title := vars["title"]
	p := cat.GetSingleCategory(id,post_title)
	categories.RenderTemplate(w,"categories", p)
}