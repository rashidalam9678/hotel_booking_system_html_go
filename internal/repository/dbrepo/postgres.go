package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/rashidalam9678/hotel_booking_system_html_go/internal/models"
	"golang.org/x/crypto/bcrypt"
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
	defer rows.Close()

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

// GetUserById return the user by id
func (m *postgresDBRepo) GetUserById(id int)(models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	stmt:=`select id, first_name, last_name, email, password, access_level, created_at, updated_at 
		where id=$1`
	row:= m.DB.QueryRowContext(ctx,stmt,id)
	err:= row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.AccessLevel,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)	
	if err!=nil{
		return user,err
	}

	return user,nil
}

func (m *postgresDBRepo) UpdateUserById(u models.User)( error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query:= `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`
	_,err:=m.DB.ExecContext(ctx,query,
			u.FirstName,
			u.LastName,
			u.Email,
			u.AccessLevel,
			time.Now(),
		)
		if err!=nil{
			return err
		}
	return 	nil	
}

//Authenticate takes the user password and retruns the user id if exist otherwise returns zero and error
func (m *postgresDBRepo) Authenticate(email,testPassword string)(int, string,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	stmt:=`select id, password from users where email=$1`
	row:= m.DB.QueryRowContext(ctx,stmt,email)

	err:=row.Scan(
		&id,
		&hashedPassword,
	)
	if err!=nil{
		return id,hashedPassword,err
	}

	err= bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(testPassword))
	if err==bcrypt.ErrMismatchedHashAndPassword{
		return 0,"",errors.New("incorrect Password")
	}else if err!=nil{
		return 0, "",err
	}

	return id,hashedPassword,nil
}

// ALl Reservation return a slice of all the reservation
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation []models.Reservation

	query:= `
			select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
			r.created_at, r.updated_at,r.processed,
			rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			order by r.start_date asc

		`
	rows,err:= m.DB.QueryContext(ctx,query,)
	if err!= nil{
		return reservation, err
	}	
	defer rows.Close()
	
	for rows.Next(){
		var i models.Reservation
		err= rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err!= nil{
			return reservation, err
		}
		reservation= append(reservation, i)
	}
	if err= rows.Err(); err!=nil{
		return reservation, err
	}

	return reservation,nil

}

// ALlNewReservation return a slice of all the reservation
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation []models.Reservation

	query:= `
			select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
			r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where processed = 0
			order by r.start_date asc

		`
	rows,err:= m.DB.QueryContext(ctx,query,)
	if err!= nil{
		return reservation, err
	}	
	defer rows.Close()
	
	for rows.Next(){
		var i models.Reservation
		err= rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err!= nil{
			return reservation, err
		}
		reservation= append(reservation, i)
	}
	if err= rows.Err(); err!=nil{
		return reservation, err
	}

	return reservation,nil

}


//getRoomById takes the roomId and returns the room model
func ( m *postgresDBRepo) GetReservationById(id int)(models.Reservation, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservation models.Reservation
	
	stmt:=`
			select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name 
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where r.id= $1
			`

	row:= m.DB.QueryRowContext(ctx,stmt,id)
	err:= row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
		
	)
	if err!= nil{
		return reservation,err
	}
	return reservation,nil
}


func (m *postgresDBRepo) UpdateReservation(u models.Reservation)( error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query:= `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5 where id=$6`
	_,err:=m.DB.ExecContext(ctx,query,
			u.FirstName,
			u.LastName,
			u.Email,
			u.Phone,
			time.Now(),
			u.ID,
		)
		if err!=nil{
			return err
		}
	return 	nil	
}

func (m *postgresDBRepo) DeleteReservationById(id int)( error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query:= `delete from reservations where id=$1`
	_,err:=m.DB.ExecContext(ctx,query,id)
		if err!=nil{
			return err
		}
	return 	nil	
}

func (m *postgresDBRepo) UpdateProcessed(id,processed int)( error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query:= `update reservations set processed=$1 where id=$2`
	_,err:=m.DB.ExecContext(ctx,query,processed,id)
		if err!=nil{
			return err
		}
	return 	nil
}




