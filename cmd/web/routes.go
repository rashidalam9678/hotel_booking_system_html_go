package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler{
	mux:= chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)


	mux.Get("/",handlers.Repo.Home)
	mux.Get("/about",handlers.Repo.About)

	fileServer:= http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))

	return mux
}