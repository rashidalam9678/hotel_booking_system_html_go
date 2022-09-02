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
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/driver"
)
const PortNumber= ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
// var warningLog *log.Logger



func main(){
	// what kind of data we will store in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

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

	



	log.Println("starting connection to database")
	db,err:= driver.ConnectSQL("host=localhost port=5432 dbname=hotel_bookings user=mr.mra password=")
	if err != nil{
		log.Fatal("can not connect to database. Dying.....")
	}
	defer db.SQL.Close()

	mailChan:= make(chan models.MailData)
	app.MailChan=mailChan

	defer close(app.MailChan)
	ListenForMail()

	fmt.Println("starting mail listener")

	


	repo:= handlers.NewRepo(&app,db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	log.Println("succesfully connected to data base")

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