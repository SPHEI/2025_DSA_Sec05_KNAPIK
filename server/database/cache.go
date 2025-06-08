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
        user_id INTEGER NOT NULL,
        token TEXT NOT NULL UNIQUE
    )`)
	return err
}

func InsertToken(db *sql.DB, userId int, token string) error {
	_, err := db.Exec(`
        INSERT INTO sessions (user_id, token)
        VALUES (?, ?)`,
		userId, token)
	return err
}

func GetToken(db *sql.DB, token string) (int, error) {
	var userId int

	err := db.QueryRow(`
        SELECT id FROM sessions 
        WHERE token = ?`, token).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("invalid session")
		}
		return -1, err
	}

	return userId, err
}

func DeleteToken(db *sql.DB, token string) {
	_, _ = db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
