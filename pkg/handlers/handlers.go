package handlers

import (
	"net/http"

	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/models"
	"github.com/rashidalam9678/hotel_booking_system_html_go/pkg/render"
)

// Repo the repository used by handlers
var Repo *Repository

// Reepository is the type of repository
type Repository struct{
	App *config.AppConfig
}

// NewRepo creates the new repository 
func NewRepo(a *config.AppConfig) *Repository{
	return &Repository{
		App : a,
	}
}

//NewHandlers sets the repo for handlers
func NewHandlers(r *Repository){
	Repo=r
}

// Home page Handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request){
		remoteIp:= r.RemoteAddr
		m.App.Session.Put(r.Context(),"remote_ip",remoteIp)

	render.RenderTemplate(w,"home.page.tmpl",&models.TemplateData{
	})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request){

	stringmap := make(map[string]string)
	stringmap["test"]= "hello from myside"
	stringmap["remote_ip"]=m.App.Session.GetString(r.Context(),"remote_ip")

	render.RenderTemplate(w, "about.page.tmpl",&models.TemplateData{
		StringMap:stringmap,
	})
}

