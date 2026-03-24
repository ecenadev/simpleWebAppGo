package main

import (
	"encoding/gob"
	"net/http"
	"time"
	"web3/pckge/config"
	"web3/pckge/handlers"
	"web3/pckge/helpers"
	"web3/pckge/models"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var app config.AppConfig

func main() {
	gob.Register(models.Article{})

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false //set to false because we dont have a https server
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = sessionManager

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()

	if err != nil {
		helpers.ErrorCheck(err)
	}

	// // http.HandleFunc("/", handlers.HomeHandler)
	// // http.HandleFunc("/about", handlers.AboutHandler)

	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	helpers.ErrorCheck(err)
	// }
}
