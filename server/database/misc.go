package database

import (
	"database/sql"
	"log"
)

func GetPassword(db *sql.DB, email string) (string, error) {
	query := "SELECT password FROM User WHERE email = ?"

	var password string
	err := db.QueryRow(query, email).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return "", nil
		}
		log.Println("Error retrieving user:", err)
		return "", err
	}

	return password, nil
}

func GetRole(db *sql.DB, email string) (int, error) {
	query := "SELECT role_id FROM User WHERE email = ?"

	var roleId int
	err := db.QueryRow(query, email).Scan(&roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, err
	}

	return roleId, nil
}

func GetId(db *sql.DB, email string) (int, error) {
	query := "SELECT id FROM User WHERE email = ?"
	var id int

	err := db.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, err
	}

	return id, nil
}
