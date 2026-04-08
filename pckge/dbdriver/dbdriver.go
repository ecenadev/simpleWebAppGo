package dbdriver

import (
	"database/sql"

	"time"

	_ "github.com/jackc/pgx/v5"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConnect = &DB{}

const maxOpenDbConns = 20
const maxIdleDbConns = 10
const maxDBLifetime = 5 * time.Minute

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectSQL(dns string) (*DB, error) {
	db, err := NewDatabase(dns)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenDbConns)
	db.SetMaxIdleConns(maxIdleDbConns)
	db.SetConnMaxLifetime(maxDBLifetime)

	dbConnect.SQL = db

	return dbConnect, nil
}
