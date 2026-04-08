package dbrepo

import (
	"database/sql"
	"web3/pckge/config"
	"web3/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}