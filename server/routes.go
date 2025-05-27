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
	router.HandleFunc("GET /tenant/list", app.getTenants)
	router.HandleFunc("GET /subcon/info", app.subInfo)

	router.HandleFunc("GET /owner/list", app.getOwners)
	router.HandleFunc("POST /owner/add", app.addOwner)

	router.HandleFunc("GET /apartament/list", app.getApartaments)
	router.HandleFunc("POST /apartament/add", app.addApartament)

	router.HandleFunc("GET /renting/current", app.getCurrentRenting)
	router.HandleFunc("POST /renting/start", app.addNewRenting)
	router.HandleFunc("POST /renting/end", app.setEndOfRenting)

	router.HandleFunc("GET /faults/list", app.getReports)
	router.HandleFunc("POST /faults/add", app.addFault)

	router.HandleFunc("GET /subcon/list", app.getContractors)
	router.HandleFunc("POST /subcon/add", app.addContractor)

	////

	router.HandleFunc("POST /adduser", app.addUser)

	router.HandleFunc("POST /subspec", app.getSubContractorSpec)
	router.HandleFunc("POST /addsubspec", app.addSubContractorSpec)
	router.HandleFunc("POST /addapartament", app.addApartament)
	router.HandleFunc("POST /changerent", app.changeRent)

	return router
}
