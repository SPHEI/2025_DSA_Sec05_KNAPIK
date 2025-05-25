package database

import (
	"database/sql"
	"log"
)

func SetupDatabase(db *sql.DB) error {
	return nil
}

func GetTenent(db *sql.DB) ([]string, error) {
	query := `SELECT User.name FROM User INNER JOIN Role ON User.role_id = Role.id WHERE Role.name = "tenent"`

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Println(err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return names, nil
}

func GetApartaments(db *sql.DB) ([]string, error) {
	query := `SELECT name FROM Apartament`

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Println(err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return names, nil
}

func AddUser(db *sql.DB, user []string, role_id int) error {
	query := "INSERT INTO User (name, password, email, phone, role_id) VALUES(?, ?, ?, ?, ?)"

	//id := getUserCount(db) + 1

	_, err := db.Exec(query, user[0], user[1], user[2], user[3], role_id)
	return err
}

func GetUser(db *sql.DB, email string) (string, error) {
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

func GetInfo(db *sql.DB, email string) (string, string, int, error) {
	query := "SELECT name, phone, role_id FROM User WHERE email = ?"

	var name, phone string
	var role_id int
	err := db.QueryRow(query, email).Scan(&name, &phone, &role_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return "", "", -1, nil
		}
		log.Println("Error retrieving user:", err)
		return "", "", -1, err
	}

	return name, phone, role_id, nil
}

func GetRole(db *sql.DB, email string) (int, error) {
	query := "SELECT role_id FROM User WHERE email = ?"

	var role_id int
	err := db.QueryRow(query, email).Scan(&role_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, err
	}

	return role_id, nil
}
