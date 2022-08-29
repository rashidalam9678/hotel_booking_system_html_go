package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/handlers"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/helpers"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/render"
)
const PortNumber= ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var warningLog *log.Logger



func main(){
	// what kind of data we will store in session
	gob.Register(models.Reservation{})

	app.InProduction=false
	session= scs.New()
	session.Lifetime= 24*time.Hour
	session.Cookie.Secure=app.InProduction
	session.Cookie.Persist= true
	session.Cookie.SameSite=http.SameSiteLaxMode

	app.Session= session

	infoLog= log.New(os.Stdout,"INFO\t", log.Ldate|log.Ltime)
	app.InfoLog= infoLog

	errorLog = log.New(os.Stdout,"ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog= errorLog


	

	
	tc,err:= render.CreateTemplateCache()
	if err!=nil {
		log.Fatal(err)
	}
	app.TemplateCache=tc
	app.UseCache= false

	repo:= handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)


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