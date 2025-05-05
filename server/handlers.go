package main

import (
	"encoding/json"
	"log"
	"net/http"

	"server/database"

	_ "github.com/glebarez/go-sqlite"
)

func prepareResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
}

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorType  string `json:"error_type"`
	// to be expanded
}

// interface for json needed
func sendError(w http.ResponseWriter, error Error) {
	if err := json.NewEncoder(w).Encode(error); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) getTenantList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	var err error
	tenants := struct {
		Names []string `json:"names"`
	}{}

	tenants.Names, err = database.GetTenent(app.DB)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	if err := json.NewEncoder(w).Encode(tenants); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
