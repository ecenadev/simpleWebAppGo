package main

import (
	"net/http"
	"web3/pckge/config"
	"web3/pckge/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// MUX http request multiplexer
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(LogRequestInfo)
	mux.Use(NoSurf)
	mux.Use(SetupSession)
	mux.Get("/", handlers.Repo.HomeHandler)
	mux.Get("/about", handlers.Repo.AboutHandler)
	mux.Get("/login", handlers.Repo.LoginHandler)
	mux.Get("/makepost", handlers.Repo.MakePostHandler)
	mux.Get("/page", handlers.Repo.PageHandler)

	mux.Post("/makepost", handlers.Repo.PostMakePostHandler)

	mux.Get("/article-recieved", handlers.Repo.ArticleReceived)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
