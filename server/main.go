package main

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
)

type app struct {
	DB    *sql.DB
	CACHE *sql.DB
}

func main() {
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	app := app{
		DB: db,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
