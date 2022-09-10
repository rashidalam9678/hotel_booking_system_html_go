package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/driver"
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
// var warningLog *log.Logger



func main(){
	db,err:= run()
	if err!= nil{
		log.Fatal(err)
	}

	defer db.SQL.Close()
	defer close(app.MailChan)

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


func run() (*driver.DB ,error){

	// what kind of data we will store in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	UseCache := flag.Bool("cache",true,"Use Template Cache")
	InProduction :=flag.Bool("production",true,"Application is in production")
	dbName :=flag.String("dbname","","database name")
	dbHost :=flag.String("dbhost","localhost","database host")
	dbUser :=flag.String("dbuser","","database user")
	dbPassword :=flag.String("dbpassword","","database password")
	dbPort :=flag.String("dbport","5432","database port")
	dbSsl :=flag.String("dbssl","disabled","database ssl setting (disable, prefer, require) ")

	flag.Parse()

	app.InProduction= *InProduction

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
	app.UseCache= *UseCache
	// app.UseCache= true

	

	connectionString:= fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",*dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSsl)

	log.Println("starting connection to database")
	// db,err:= driver.ConnectSQL("host=localhost port=5432 dbname=hotel_bookings user=mr.mra password=")
	db,err:= driver.ConnectSQL(connectionString)
	if err != nil{
		log.Fatal("can not connect to database. Dying.....")
	}
	

	mailChan:= make(chan models.MailData)
	app.MailChan=mailChan

	ListenForMail()

	fmt.Println("starting mail listener")

	


	repo:= handlers.NewRepo(&app,db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	log.Println("succesfully connected to data base")

	return db, nil

}