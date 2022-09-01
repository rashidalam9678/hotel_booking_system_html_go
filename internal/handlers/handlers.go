package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/driver"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/forms"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/helpers"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/render"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/repository"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.Database
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:dbrepo.NewPostgresRepo(db.SQL,a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	sd:=r.Form.Get("start_date")
	ed:=r.Form.Get("end_date")
	layout:="2006-01-02"
	startDate,err:=time.Parse(layout,sd)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	endDate,err:=time.Parse(layout,ed)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	roomId,err:= strconv.Atoi(r.Form.Get("room_id"))
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationId,err:= m.DB.InsertReservation(reservation)

	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	restriction:= models.RoomRestriction{
		StartDate :startDate,
		EndDate: endDate,
		RoomId: roomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}

	err=m.DB.InsertRoomRestriction(restriction)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	m.App.Session.Put(r.Context(),"reservation",reservation)
	http.Redirect(w,r,"/reservation-summary",http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w,err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary Renders the reservation summary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation,ok:= m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		m.App.ErrorLog.Println("can't get information from session")
		m.App.Session.Put(r.Context(),"error","can't get reservation from session")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(),"rservation")
	data:= make(map[string]interface{})
	data["reservation"]=reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:data,
	})
}