package models

import "time"



// User struct hold the user model
type User struct {
	ID int
	FirstName string
	LastName string
	Email string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
	AccessLevel int
}

type Room struct{
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reservation struct{
	ID int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomId int
	CreatedAt time.Time
	UpdatedAt time.Time

}

type RoomRestriction struct{
	ID int
	StartDate time.Time
	EndDate time.Time
	RoomId int
	ReservationId int
	CreatedAt time.Time
	UpdatedAt time.Time
	RestrictionId int
	Reservatin Reservation
	Room Room
}

type Restriction struct{
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}