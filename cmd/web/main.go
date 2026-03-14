package main

import (
	"net/http"
	handlers "web3/pckge/handlers"
	helpers "web3/pckge/helpers"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/about", handlers.AboutHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		helpers.ErrorCheck(err)
	}
}
