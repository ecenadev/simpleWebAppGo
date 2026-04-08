package helpers

import (
	"log"
	"net/http"
	"web3/pckge/config"
)

var app config.AppConfig

func ErrorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsAuthenticated(r *http.Request) bool {
	exist := app.Session.Exists(r.Context(), "user_id")
	// = app.Session.Remove(r.Context(), "user_id")
	return exist
}
