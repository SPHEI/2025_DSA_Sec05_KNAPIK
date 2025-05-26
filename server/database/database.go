package database

import (
	"database/sql"
	"log"
)

func GetEmails(db *sql.DB) ([]string, error) {
	query := `SELECT User.email FROM User`

	data, err := getMultiRow(db, query)
	return data, err

}

func GetTenent(db *sql.DB) ([]string, error) {
	query := `SELECT User.name FROM User INNER JOIN Role ON User.role_id = Role.id WHERE Role.name = "tenent"`

	data, err := getMultiRow(db, query)

	return data, err
}

func GetSubcontractorSpec(db *sql.DB) ([]string, error) {
	query := `SELECT Speciality.name FROM Speciality`

	data, err := getMultiRow(db, query)

	return data, err
}

func AddSpec(db *sql.DB, name string) error {
	query := "INSERT INTO Speciality (name) VALUES(?)"
	err := insertValue(db, query, name)
	return err
}

func GetApartaments(db *sql.DB) ([]string, error) {
	query := `SELECT name FROM Apartament`

	data, err := getMultiRow(db, query)
	return data, err
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
	var roleId int
	err := db.QueryRow(query, email).Scan(&name, &phone, &roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return "", "", -1, nil
		}
		log.Println("Error retrieving user:", err)
		return "", "", -1, err
	}

	return name, phone, roleId, nil
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

func ChangeRent(db *sql.DB, apartamentId int, rent string) error {
	query := "INSERT INTO Pricing_History (apartament_id, rent) VALUES(?, ?)"

	//id := getUserCount(db) + 1

	_, err := db.Exec(query, apartamentId, rent)
	return err
}
