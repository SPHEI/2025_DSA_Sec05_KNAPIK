package database

import (
	"database/sql"
	"log"
)

func GetInfo(db *sql.DB, email string) (int, string, string, int, error) {
	query := "SELECT id, name, phone, role_id FROM User WHERE email = ?"
	var name, phone string
	var id, roleId int

	err := db.QueryRow(query, email).Scan(&id, &name, &phone, &roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, "", "", -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, "", "", -1, err
	}

	return id, name, phone, roleId, nil
}

func GetApartamentId(db *sql.DB, userId int) (int, error) {
	query := `SELECT apartment_id 
	FROM Renting_History WHERE end_date IS NULL AND user_id = ?`
	var apartamentId int

	err := db.QueryRow(query, userId).Scan(&apartamentId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, err
	}

	return apartamentId, err
}

func GetSubconInfo(db *sql.DB, userId int) (string, string, int, error) {
	query := `SELECT address, NIP, speciality_id FROM Subcontractor WHERE user_id = ?`
	var address, NIP string
	var speciality_id int

	err := db.QueryRow(query, userId).Scan(&address, &NIP, &speciality_id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return "", "", -1, nil
		}
		log.Println("Error retrieving user:", err)
		return "", "", -1, err
	}

	return address, NIP, speciality_id, err

}

func GetApartamentsData(db *sql.DB) ([]int, []string, []string, []string, []string, []string, []int, error) {
	query := `SELECT id FROM Apartament`
	id, err := getMultiRowInt(db, query)

	query = `SELECT name FROM Apartament`
	name, err := getMultiRow(db, query)

	query = `SELECT street FROM Apartament`
	street, err := getMultiRow(db, query)

	query = `SELECT building_number FROM Apartament`
	buildingNumber, err := getMultiRow(db, query)

	query = `SELECT building_name FROM Apartament`
	buildingName, err := getMultiRow(db, query)

	query = `SELECT flat_number FROM Apartament`
	flatNumber, err := getMultiRow(db, query)

	query = `SELECT owner_id FROM Apartament`
	ownerId, err := getMultiRowInt(db, query)

	return id, name, street, buildingNumber, buildingName, flatNumber, ownerId, err
}

func AddApartament(db *sql.DB, data []string, ownerId int) error {
	query := "INSERT INTO Apartament (name, street, building_number, building_name, flat_number, owner_id) VALUES(?, ?, ?, ?, ?, ?)"

	_, err := db.Exec(query, data[0], data[1], data[2], data[3], data[4], ownerId)

	return err
}

func GetOwners(db *sql.DB) ([]int, []string, []string, []string, error) {
	query := `SELECT id FROM Owner`
	id, err := getMultiRowInt(db, query)

	query = `SELECT name FROM Owner`
	name, err := getMultiRow(db, query)

	query = `SELECT email FROM Owner`
	street, err := getMultiRow(db, query)

	query = `SELECT phone FROM Owner`
	phone, err := getMultiRow(db, query)

	return id, name, street, phone, err
}

func AddOwner(db *sql.DB, data []string) error {
	query := "INSERT INTO Owner (name, email, phone) VALUES(?, ?, ?)"

	_, err := db.Exec(query, data[0], data[1], data[2])

	return err
}

func GetRent(db *sql.DB, apartamentId int) (float32, error) {
	query := `SELECT price FROM Pricing_History WHERE is_current = 0 AND apartment_id = ?`
	var rent float32

	err := db.QueryRow(query, apartamentId).Scan(&rent)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, err
	}

	return rent, err

}

//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////

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

func AddUser(db *sql.DB, user []string, role_id int) error {
	query := "INSERT INTO User (name, password, email, phone, role_id) VALUES(?, ?, ?, ?, ?)"

	//id := getUserCount(db) + 1

	_, err := db.Exec(query, user[0], user[1], user[2], user[3], role_id)
	return err
}

func ChangeRent(db *sql.DB, apartamentId int, rent float32) error {
	query := "INSERT INTO Pricing_History (apartament_id, rent) VALUES(?, ?)"

	//id := getUserCount(db) + 1

	_, err := db.Exec(query, apartamentId, rent)
	return err
}
