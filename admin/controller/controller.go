package controller

import (
	"html/template"
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/config"
	"fmt"
)

/* -- RenderTemplate --
 * 	The function template.ParseFiles will read the contents of multiple "name".html files into cache.
 *	The method t.Execute executes the template, the string must correspond to the name giving to the template
 *	when defining them.
 *	After executing all the subtemplates, t.Execute will write the generated HTML to the http.ResponseWriter.
 *  a declared empty interface can take in any type:
 *  Ik kan de Data interface een method geven die elk Type dat wordt gevoerd aan de
 * rendertemplate functie moet hebben als voorwaarde. Iets zoals PHP abstract class.
 * http://go-book.appspot.com/interfaces.html
*/

type Data interface {
	// GetPost() error
}

// Renders the specified files/template and can take in any type, the template file will define if it can use the given type.
func RenderTemplate(w http.ResponseWriter, name string, data interface{}){
	fmt.Println(data)
	t, err := template.ParseFiles(config.Templates+"/"+"header.html",config.Templates+"/"+"nav.html",config.View + "/" + name + ".html",config.Templates+"/"+"footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w,"header",nil)
	t.ExecuteTemplate(w,"nav",nil)
	t.ExecuteTemplate(w,name,data)
	t.ExecuteTemplate(w,"footer",nil)
	err = t.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
