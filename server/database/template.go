package database

import (
	"database/sql"
	"log"
)

func insertValue(db *sql.DB, query, value string) error {
	_, err := db.Exec(query, value)
	log.Println(err)
	return err
}

func getMultiRow(db *sql.DB, query string) ([]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			//
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func getMultiRowInt(db *sql.DB, query string) ([]int, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []int
	for rows.Next() {
		var value int
		if err := rows.Scan(&value); err != nil {
			//
		}
		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}
