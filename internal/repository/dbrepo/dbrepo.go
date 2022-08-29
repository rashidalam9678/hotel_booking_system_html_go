package dbrepo

import (
	"database/sql"

	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/config"
	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/repository"
)

type postgresDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig ) repository.Database{
	return &postgresDBRepo{
		App:a,
		DB:conn,
	}
}