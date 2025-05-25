package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)
	router.HandleFunc("POST /adduser", app.addUser)
	router.HandleFunc("GET /info", app.info)
	router.HandleFunc("GET /tenents", app.getTenantList)
	router.HandleFunc("GET /apartaments", app.getApartamentList)

	return router
}
