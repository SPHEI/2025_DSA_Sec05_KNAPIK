package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net/http"
	"server/database"
	"server/sqlc"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed schema.sql
var ddl string

type app struct {
	DB    *sql.DB
	CACHE *sql.DB
	Query *sqlc.Queries
	Ctx   context.Context
}

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	cache, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(ctx, ddl)
	if err != err {
		log.Println(err)
	}

	if err = database.SetupCache(cache); err != nil {
		log.Fatal(err)
	}

	app := app{
		DB:    db,
		CACHE: cache,
		Query: sqlc.New(db),
		Ctx:   ctx,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
