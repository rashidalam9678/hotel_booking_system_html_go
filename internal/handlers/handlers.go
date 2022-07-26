package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok{
		helpers.ServerError(w,errors.New("can not get reservation from session"))
	}

	room,err:=m.DB.GetRoomById(res.RoomId)
	if err!= nil{
		helpers.ServerError(w,err)
	}
	res.Room.RoomName=room.RoomName

	m.App.Session.Put(r.Context(),"reservation",res)

	sd:=res.StartDate.Format("2006-01-02")
	ed:=res.EndDate.Format("2006-01-02")

	stringMap:= make(map[string]string)
	stringMap["start_date"]=sd
	stringMap["end_date"]=ed

	if !ok {
		helpers.ServerError(w,errors.New("can't get reservation details from session"))
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation,ok:=m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		helpers.ServerError(w,errors.New("can't get the reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	// sd:=r.Form.Get("start_date")
	// ed:=r.Form.Get("end_date")
	// layout:="2006-01-02"
	// startDate,err:=time.Parse(layout,sd)
	// if err!= nil{
	// 	helpers.ServerError(w,err)
	// 	return
	// }

	// endDate,err:=time.Parse(layout,ed)
	// if err!= nil{
	// 	helpers.ServerError(w,err)
	// 	return
	// }

	// roomId,err:= strconv.Atoi(r.Form.Get("room_id"))

	// if err!= nil{
	// 	helpers.ServerError(w,err)
	// 	return
	// }

	
	reservation.FirstName=r.Form.Get("first_name")
	reservation.LastName=r.Form.Get("last_name")
	reservation.Email=r.Form.Get("email")
	reservation.Phone=r.Form.Get("phone")


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
		StartDate :reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomId: reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}

	err=m.DB.InsertRoomRestriction(restriction)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	m.App.Session.Put(r.Context(),"reservation",reservation)

	html:=fmt.Sprintf(`
				<strong> Reservation Notification</strong><br>
				A reservation has been made for %s from %s to %s.
			`,reservation.FirstName,reservation.StartDate.Format("2006-02-01"),reservation.EndDate.Format("2006-02-01"))

	msg:=models.MailData{
		To:reservation.Email,
		From:"me@here.com",
		Subject: "Reservation Confirmation",
		Content: html,
		Template: "basic.html",
	}
	m.App.MailChan <- msg


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
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{
	})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
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

	availableRooms,err:=m.DB.SearchAvailablityForAllRooms(startDate,endDate)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}
	for _,j := range availableRooms{
		m.App.InfoLog.Println("ROOMS",j.ID, j.RoomName)
	}
	if len(availableRooms)==0{
		
		m.App.Session.Put(r.Context(),"error","No Rooms Available")
		m.App.InfoLog.Println("No rooms are availble")
		http.Redirect(w,r,"search-availability",http.StatusSeeOther)
		
	}
	data:= make(map[string]interface{})
	data["rooms"]=availableRooms

	m.App.Session.Put(r.Context(),"reservation",models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	})

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data:data,
	})
	
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId,err:= strconv.Atoi(chi.URLParam(r,"id"))
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w,err)
		return
	}

	res.RoomId=roomId
	m.App.Session.Put(r.Context(),"reservation",res)

	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)


	
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request){
	roomId,err:= strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	sd:=r.URL.Query().Get("s")
	ed:=r.URL.Query().Get("e")

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

	var res models.Reservation

	res.RoomId=roomId
	res.StartDate=startDate
	res.EndDate= endDate

	room,err:=m.DB.GetRoomById(roomId)
	if err!= nil{
		helpers.ServerError(w,err)
	}
	res.Room.RoomName=room.RoomName

	m.App.Session.Put(r.Context(),"reservation",res)

	http.Redirect(w,r,"make-reservation",http.StatusSeeOther)


}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	RoomId string   `json:"room_id"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd:= r.Form.Get("start")
	ed:= r.Form.Get("end")

	layout:="2006-01-02"
	startDate,err:=time.Parse(layout,sd)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	endDate,err:=time.Parse(layout,ed)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	roomId,err:=strconv.Atoi(r.Form.Get("room_id"))
	if err!=nil{
		helpers.ServerError(w,err)
		return
	}

	available,err:= m.DB.SearchAvailablityByDatesByRoomId(startDate,endDate,roomId)
	if err!=nil{
		helpers.ServerError(w,err)
		return
	}

	resp := jsonResponse{
		OK:      available,
		Message: "",
		StartDate:sd,
		EndDate: ed,
		RoomId:strconv.Itoa(roomId) ,
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
	sd:= reservation.StartDate.Format("2006-01-02")
	ed:= reservation.EndDate.Format("2006-01-02")

	stringMap:= make(map[string]string)
	stringMap["start_date"]=sd
	stringMap["end_date"]=ed

	m.App.Session.Remove(r.Context(),"rservation")
	data:= make(map[string]interface{})
	data["reservation"]=reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:data,
		StringMap: stringMap,
	})
}


// ShowLogin displays the login page
func (m *Repository)ShowLogin(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}
// ShowSignUp displays the login page
func (m *Repository)ShowSignup(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "signup.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowSignup( w http.ResponseWriter, r *http.Request) {
	_= m.App.Session.RenewToken(r.Context())
	err:= r.ParseForm()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	form:= forms.New(r.Form)

	// email := r.Form.Get("email")
	// password:= r.Form.Get("password")


	form.Required("email","password")
	form.IsEmail("email")

	if !form.Valid(){
		//take user back to page
		render.Template(w,r,"login.page.tmpl",&models.TemplateData{
			Form:form,
		})
		return
	}


}



func (m *Repository)PostShowLogin(w http.ResponseWriter, r *http.Request){
	_= m.App.Session.RenewToken(r.Context())

	err:= r.ParseForm()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	form:= forms.New(r.Form)

	email := r.Form.Get("email")
	password:= r.Form.Get("password")

	form.Required("email","password")
	form.IsEmail("email")

	if !form.Valid(){
		//take user back to page
		render.Template(w,r,"login.page.tmpl",&models.TemplateData{
			Form:form,
		})
		return
	}

	 id,_,err:=m.DB.Authenticate(email,password)
	 if err!= nil{
		log.Println(err)
		m.App.Session.Put(r.Context(),"error","Invalid credentials")
		http.Redirect(w,r,"/user/login",http.StatusSeeOther)
		return
	 }

	m.App.Session.Put(r.Context(),"user_id",id)
	m.App.Session.Put(r.Context(),"success","Logged in successfully")
	http.Redirect(w,r,"/",http.StatusSeeOther)

}


func (m *Repository) Logout(w http.ResponseWriter, r *http.Request){
	err:=m.App.Session.Destroy(r.Context())
	if err !=nil{
		helpers.ServerError(w,err)
	}
	err=m.App.Session.RenewToken(r.Context())
	if err !=nil{
		helpers.ServerError(w,err)
	}
	http.Redirect(w,r,"/user/login",http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request){
	render.Template(w,r,"admin-dashboard.page.tmpl",&models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request){

	reservations, err:= m.DB.AllNewReservations()
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}
	data:= make(map[string] interface{})
	data["reservations"]=reservations

	render.Template(w,r,"admin-reservations-new.page.tmpl",&models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request){
	reservations, err:= m.DB.AllReservations()
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}
	data:= make(map[string] interface{})
	data["reservations"]=reservations

	render.Template(w,r,"admin-reservations-all.page.tmpl",&models.TemplateData{
		Data: data,
	})
}


func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request){
	explode := strings.Split(r.RequestURI,"/")
	id,err:= strconv.Atoi(explode[4])
	if err!=nil{
		helpers.ServerError(w,err)
	}


	src:= explode[3]

	stringMap:= make(map[string] string)
	stringMap["src"]=src


	reservation, err:= m.DB.GetReservationById(id)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}
	data:= make(map[string] interface{})
	data["reservation"]=reservation

	render.Template(w,r,"admin-reservations-show.page.tmpl",&models.TemplateData{
		Data: data,
		StringMap: stringMap,
		Form: forms.New(nil),
	})
}

func (m *Repository) AdminUpdateReservation(w http.ResponseWriter, r *http.Request){

	explode := strings.Split(r.RequestURI,"/")
	id,err:= strconv.Atoi(explode[4])
	if err!=nil{
		helpers.ServerError(w,err)
		return
	}

	src:= explode[3]

	err= r.ParseForm()
	if err!=nil{
		helpers.ServerError(w,err)
		return
	}

	reservation, err:= m.DB.GetReservationById(id)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	reservation.FirstName=r.Form.Get("first_name")
	reservation.LastName=r.Form.Get("last_name")
	reservation.Email=r.Form.Get("email")
	reservation.Phone=r.Form.Get("phone")

	err= m.DB.UpdateReservation(reservation)
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	m.App.Session.Put(r.Context(),"flash","reservation Updated Successfully !")
	http.Redirect(w,r,fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)

}



func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request){
	id,err:= strconv.Atoi(chi.URLParam(r,"id"))
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}

	src:=chi.URLParam(r,"src")

	err= m.DB.UpdateProcessed(id,1)
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}
	m.App.Session.Put(r.Context(),"flash","Reservation Marked as processed")
	http.Redirect(w,r,fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)

}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request){
	id,err:= strconv.Atoi(chi.URLParam(r,"id"))
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}

	src:=chi.URLParam(r,"src")

	err= m.DB.DeleteReservationById(id)
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}
	m.App.Session.Put(r.Context(),"flash","Reservation Deleted ")
	http.Redirect(w,r,fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)

}


func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request){
	render.Template(w,r,"admin-reservations-calendar.page.tmpl",&models.TemplateData{})
}

