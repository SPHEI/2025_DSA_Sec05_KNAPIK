package database

import (
	"database/sql"
	"log"
)

func SetupDatabase(db *sql.DB) error {
	return nil
}

func GetEmails(db *sql.DB) ([]string, error) {
	query := `SELECT User.email FROM User`

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

func GetSubcontractorSpec(db *sql.DB, subContractor_email string) ([]string, error) {
	query := `SELECT Speciality.name FROM User 
	INNER JOIN Subcontractor 
	ON Subcontractor.user_id = User.id 
	INNER JOIN Speciality 
	ON Subcontractor.speciality_id = Speciality.id 
	WHERE User.email = ?`

	rows, err := db.Query(query, subContractor_email)
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

func AddApartament(db *sql.DB, data []string) error {
	query := "SELECT id FROM Owner WHERE email = ?"

	var owner_id int
	err := db.QueryRow(query, data[0]).Scan(&owner_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return nil
		}
		log.Println("Error retrieving user:", err)
		return err
	}

	query = "INSERT INTO Apartament (name, street, building_namber, building_name, flat_number, owner_id) VALUES(?, ?, ?, ?, ?)"

	_, err = db.Exec(query, data[1], data[2], data[3], data[4], data[5], owner_id)
	return err
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
