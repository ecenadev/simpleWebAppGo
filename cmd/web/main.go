package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"
	"web3/pckge/config"
	"web3/pckge/dbdriver"
	"web3/pckge/handlers"
	"web3/pckge/helpers"
	"web3/pckge/models"
	"web3/pckge/render"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var app config.AppConfig

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

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

func run() (*dbdriver.DB, error) {
	gob.Register(models.Article{})
	gob.Register(models.Post{})
	gob.Register(models.User{})

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false //set to false because we dont have a https server
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = sessionManager

	db, err := dbdriver.ConnectSQL("host=localhost port=5432 dbname= blog_db user=postgres password = ****** ")

	if err != nil {
		log.Fatal("Can't connect to database")
	}

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewAppConfig(&app)

	return db, nil
}
