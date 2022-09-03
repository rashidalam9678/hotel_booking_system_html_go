package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler{
	mux:= chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)


	mux.Get("/",handlers.Repo.Home)
	mux.Get("/about",handlers.Repo.About)	
	mux.Get("/contact",handlers.Repo.Contact)	
	mux.Get("/search-availability",handlers.Repo.Availability)
	mux.Post("/search-availability",handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/generals-quarters",handlers.Repo.Generals)
	mux.Get("/majors-suite",handlers.Repo.Majors)
	mux.Get("/choose-room/{id}",handlers.Repo.ChooseRoom)
	mux.Get("/book-room",handlers.Repo.BookRoom)
	mux.Get("/make-reservation",handlers.Repo.Reservation)
	mux.Post("/make-reservation",handlers.Repo.PostReservation)
	mux.Get("/reservation-summary",handlers.Repo.ReservationSummary)

	mux.Get("/user/login",handlers.Repo.ShowLogin)
	mux.Post("/user/login",handlers.Repo.PostShowLogin)
	mux.Get("/user/logout",handlers.Repo.Logout)

	mux.Route("/admin",func(r chi.Router) {
			r.Use(Auth)
			r.Get("/admin",handlers.Repo.AdminDashboard)
	})

	fileServer:= http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))

	return mux
}