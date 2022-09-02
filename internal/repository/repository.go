package repository

import (
	"time"

	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
)

type Database interface{
	AllUsers() bool
	InsertReservation(res models.Reservation) (int,error)
	InsertRoomRestriction( res models.RoomRestriction)(error)
	SearchAvailablityByDatesByRoomId(start time.Time,end time.Time , roomId int)(bool, error)
	SearchAvailablityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int)(models.Room, error)
}