package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/models"
)

var app *config.AppConfig
// NewTemplates sets the config for the new templates
func NewTemplates(a *config.AppConfig){
	app= a
}

// RenderTemplate is function which can parse the given template and render it to browser
func RenderTemplate(w http.ResponseWriter, tmpl string,td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache{
		tc= app.TemplateCache
	}else{
		tc,_=CreateTemplateCache()
	}
	t,ok:=tc[tmpl]
	if !ok {
		log.Fatal("couldn't get the template cache")
	}
	buf := new(bytes.Buffer)
	err:= t.Execute(buf, td)
	if err!=nil{
		log.Println(err)
	}
	_,err= buf.WriteTo(w)
	if err!=nil {
		log.Println(err)
	}
	
}


func CreateTemplateCache()(map[string]*template.Template,error){
	myCache := map[string]*template.Template{}

	pages, err:= filepath.Glob("./templates/*.page.tmpl")
	if err != nil{
		return myCache,err
	}
	for _,page := range pages{
		name := filepath.Base(page)
		ts,err := template.New(name).ParseFiles(page)
		if err != nil{
			return myCache, err
		}

		matches,err:= filepath.Glob("./templates/*.layout.tmpl")
		if err != nil{
			return myCache, err
		}
		if len(matches)>0{
			ts,err= ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil{
				return myCache, err
			}
		}
		myCache[name]=ts

	}

	return myCache,nil
}