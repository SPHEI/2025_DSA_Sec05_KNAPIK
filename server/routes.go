package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /tenents", app.getTenantList)

	return router
}
