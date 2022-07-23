package main
import (
	"net/http"
	"html/template"
	"fmt"

)

// Home page Handler
func Home(w http.ResponseWriter, r *http.Request){
	renderTemplate(w,"home.html")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request){
	renderTemplate(w, "about.html")
}

// renderTemplate is function which can parse the given template and render it to browser
func renderTemplate(w http.ResponseWriter, tmpl string) {
	ParsedTemplate, err :=template.ParseFiles("./templates/"+ tmpl)
	if err != nil{
		fmt.Println("Unable to serve html template", err)
	}
	ParsedTemplate.Execute(w,nil)
}