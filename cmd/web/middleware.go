package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

//NoSurf protects our app from csrf
func NoSurf(next http.Handler) http.Handler{
	csrfHandler:=nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure: false,
		Path: "/",
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad saves and load session on every request
func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}