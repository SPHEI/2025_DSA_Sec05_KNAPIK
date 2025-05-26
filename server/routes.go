package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)

	router.HandleFunc("GET /info", app.info)
	router.HandleFunc("GET /tenant/info", app.tenantInfo)
	router.HandleFunc("GET /subcon/info", app.subInfo)

	router.HandleFunc("GET /owner/list", app.getOwners)
	router.HandleFunc("POST /owner/add", app.addOwner)

	router.HandleFunc("GET /apartament/list", app.getApartaments)
	router.HandleFunc("POST /apartament/add", app.addApartament)

	////

	router.HandleFunc("POST /adduser", app.addUser)

	router.HandleFunc("POST /emails", app.getEmailList)
	router.HandleFunc("POST /tenents", app.getTenantList)
	router.HandleFunc("POST /subspec", app.getSubContractorSpec)
	router.HandleFunc("POST /addsubspec", app.addSubContractorSpec)
	router.HandleFunc("POST /addapartament", app.addApartament)
	router.HandleFunc("POST /changerent", app.changeRent)

	return router
}
