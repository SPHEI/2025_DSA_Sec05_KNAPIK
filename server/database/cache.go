package database

import (
	"database/sql"
	"errors"
)

func SetupCache(db *sql.DB) error {
	// Create sessions table
	_, err := db.Exec(`
    CREATE TABLE sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        login TEXT NOT NULL,
        token TEXT NOT NULL UNIQUE
    )`)
	return err
}

func InsertToken(db *sql.DB, login, token string) error {
	_, err := db.Exec(`
        INSERT INTO sessions (login, token)
        VALUES (?, ?)`,
		login, token)
	return err
}

func GetToken(db *sql.DB, token string) (string, error) {
	var login string

	err := db.QueryRow(`
        SELECT login FROM sessions 
        WHERE token = ?`, token).Scan(&login)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid session")
		}
		return "", err
	}

	return login, err
}

func DeleteToken(db *sql.DB, token string) {
	_, _ = db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
