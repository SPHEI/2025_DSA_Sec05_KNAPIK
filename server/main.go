package main

import (
	"database/sql"
	"log"
	"net/http"
	"server/database"

	_ "github.com/glebarez/go-sqlite"
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

	cache, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	database.SetupCache(cache)

	app := app{
		DB:    db,
		CACHE: cache,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
