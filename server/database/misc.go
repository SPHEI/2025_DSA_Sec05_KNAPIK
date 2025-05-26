package database

import (
	"database/sql"
	"log"
)

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
