package main

import (
	"log"
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)

	router.Handle("GET /info", logger(http.HandlerFunc(app.info)))
	router.Handle("GET /tenant/info", logger(http.HandlerFunc(app.tenantInfo)))
	router.Handle("GET /tenant/list", logger(http.HandlerFunc(app.getTenants)))
	router.Handle("GET /subcon/info", logger(http.HandlerFunc(app.subInfo)))

	router.Handle("GET /apartament/list", logger(http.HandlerFunc(app.getApartaments)))
	router.Handle("POST /apartament/add", logger(http.HandlerFunc(app.addApartament)))

	router.Handle("GET /renting/current", logger(http.HandlerFunc(app.getCurrentRenting)))
	router.Handle("POST /renting/start", logger(http.HandlerFunc(app.addNewRenting)))
	router.Handle("POST /renting/end", logger(http.HandlerFunc(app.setEndOfRenting)))

	router.Handle("GET /faults/list", logger(http.HandlerFunc(app.getReports)))
	router.Handle("POST /faults/add", logger(http.HandlerFunc(app.addFault)))
	router.Handle("POST /faults/status", logger(http.HandlerFunc(app.updateFault)))

	router.Handle("GET /subcon/list", logger(http.HandlerFunc(app.getContractors)))
	router.Handle("POST /subcon/add", logger(http.HandlerFunc(app.addContractor)))

	router.Handle("GET /repair/list", logger(http.HandlerFunc(app.getRepairs)))
	router.Handle("POST /repair/add", logger(http.HandlerFunc(app.addRepair)))
	router.Handle("POST /repair/contractor", logger(http.HandlerFunc(app.assignSubContractor)))
	router.Handle("POST /repair/data", logger(http.HandlerFunc(app.updateRepairData)))

	router.Handle("GET /payments/list", logger(http.HandlerFunc(app.getPayments)))
	router.Handle("POST /payments/pay", logger(http.HandlerFunc(app.pay)))

	////

	router.Handle("POST /adduser", logger(http.HandlerFunc(app.addUser)))

	router.Handle("POST /subspec", logger(http.HandlerFunc(app.getSubContractorSpec)))
	router.Handle("POST /addsubspec", logger(http.HandlerFunc(app.addSubContractorSpec)))
	router.Handle("POST /addapartament", logger(http.HandlerFunc(app.addApartament)))
	router.Handle("POST /changerent", logger(http.HandlerFunc(app.changeRent)))

	return router
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}
