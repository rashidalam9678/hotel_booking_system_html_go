package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/handlers"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/render"
)
const PortNumber= ":8080"
var app config.AppConfig
var session *scs.SessionManager



func main(){
	app.InProduction=false
	session= scs.New()
	session.Lifetime= 24*time.Hour
	session.Cookie.Secure=app.InProduction
	session.Cookie.Persist= true
	session.Cookie.SameSite=http.SameSiteLaxMode

	app.Session= session
	

	
	tc,err:= render.CreateTemplateCache()
	if err!=nil {
		log.Fatal(err)
	}
	app.TemplateCache=tc
	app.UseCache= false

	repo:= handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)


	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	fmt.Println("Started Serven on: http://localhost:8080 ")

	srv:= &http.Server{
		Addr: PortNumber,
		Handler:routes(&app),
	}
	err= srv.ListenAndServe()
	if err!=nil{
		log.Fatal("Unable to start the server")
	}

	
	
}