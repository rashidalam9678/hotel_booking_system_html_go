package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation insert the booking into reservation table
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int

	stmt := `insert into reservations (first_name, last_name, email, phone,start_date, end_date, room_id,created_at, updated_at) values 
		($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id	`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,created_at,
			updated_at,restriction_id) values 
			($1,$2,$3,$4,$5,$6,$7)	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		res.ReservationId,
		res.CreatedAt,
		res.UpdatedAt,
		res.RestrictionId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailablityByDatesByRoomId(start time.Time, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var noOfRows int

	stmt := `select count(id) from room_restrictions where
			room_id=$1 and $2 < end_date and $3 > start_date
		`
	err := m.DB.QueryRowContext(ctx, stmt,
		roomId, start, end,
	).Scan(&noOfRows)

	if err != nil {
		return false, err
	}

	if noOfRows == 0 {
		return true, nil
	}
	return false, nil

}

// SearchAvailabilityForAllRooms returns the slice of all rooms available
func (m *postgresDBRepo) SearchAvailablityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var availableRooms []models.Room
	query := `
			select 
			r.id, r.room_name
			from rooms r
			where 
			r.id not in (select rr.room_id from room_restrictions rr where $1< rr.end_date  and $2>rr.start_date  )
		`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return availableRooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return availableRooms, err
		}
		availableRooms = append(availableRooms, room)
	}
	if err = rows.Err(); err!=nil{
		log.Fatal("Error scanning rows", err)
		return availableRooms,err
	}

	return availableRooms, nil
}

//getRoomById takes the roomId and returns the room model

func ( m *postgresDBRepo) GetRoomById(id int)(models.Room, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room models.Room
	
	stmt:=`
			select id, room_name from rooms where id=$1
		`

	row:= m.DB.QueryRowContext(ctx,stmt,id)
	err:= row.Scan(
		&room.ID,
		&room.RoomName,
	)
	if err!= nil{
		return room,err
	}
	return room,nil
}


