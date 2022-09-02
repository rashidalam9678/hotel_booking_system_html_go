package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
)

type AppConfig struct{
	UseCache bool
	TemplateCache map[string] *template.Template
	InfoLog 	*log.Logger
	ErrorLog 	*log.Logger
	WarningLog 	*log.Logger
	InProduction bool
	Session *scs.SessionManager
	CSRFToken string
	MailChan chan models.MailData
}