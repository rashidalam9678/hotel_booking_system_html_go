package driver

import (
	"database/sql"
	"time"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct{
	SQL *sql.DB
}

var dbCon= &DB{}

const maxOpenDbConn=10
const maxIdleDbConn=5
const maxDbLifeTime= 5*time.Minute

// ConnectSQL creates database pool for postgres
func connectSQL(dns string)(*DB ,error){
	db, err:= newDatabase(dns)
	if err !=nil{
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	dbCon.SQL=db
	err=testDB(db)

	if err !=nil{
		return nil,err
	}

	return dbCon,nil

}

// testDB tries to ping the database
func testDB(d *sql.DB) error{
	err:= d.Ping()
	if err!=nil{
		return err
	}
	return nil
}

//newDatabase creates the new database for the connection pool
func newDatabase(dns string) (*sql.DB,error){
	db,err:= sql.Open("pgx",dns)
	if err !=nil{
		return nil,err
	}
	err=db.Ping()
	if err!= nil{
		return nil,err
	}
	return db,nil
}