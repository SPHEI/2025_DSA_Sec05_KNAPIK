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
	router.HandleFunc("POST /emails", app.getEmailList)
	router.HandleFunc("POST /tenents", app.getTenantList)
	router.HandleFunc("POST /apartaments", app.getApartamentList)
	router.HandleFunc("POST /subspec", app.getSubContractorSpec)
	router.HandleFunc("POST /addsubspec", app.addSubContractorSpec)
	router.HandleFunc("POST /addapartament", app.addApartament)

	return router
}
