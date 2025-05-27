package database

import (
	"database/sql"
	"log"
	"time"
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

func GetContractorsData(db *sql.DB) ([]int, []int, []string, []string, []int, error) {
	query := `SELECT id FROM Subcontractor`
	id, err := getMultiRowInt(db, query)

	query = `SELECT user_id FROM Subcontractor`
	userId, err := getMultiRowInt(db, query)

	query = `SELECT address FROM Subcontractor`
	address, err := getMultiRow(db, query)

	query = `SELECT NIP FROM Subcontractor`
	nip, err := getMultiRow(db, query)

	query = `SELECT speciality_id FROM Subcontractor`
	specialityId, err := getMultiRowInt(db, query)

	return id, userId, address, nip, specialityId, err
}

func AddContractor(db *sql.DB, data []string, userId, specialityId int) error {
	query := "INSERT INTO Subcontractor (user_id, address, NIP, speciality_id) VALUES(?, ?, ?, ?)"

	_, err := db.Exec(query, userId, data[0], data[1], specialityId)

	return err
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

func ChangeRent(db *sql.DB, apartamentId int, rent float32) error {
	query := `UPDATE Pricing_History
	SET is_current = 1
	WHERE is_current = 0 AND apartment_id = ?`
	_, err := db.Exec(query, apartamentId)

	query = "INSERT INTO Pricing_History (apartment_id, price) VALUES(?, ?)"
	_, err = db.Exec(query, apartamentId, rent)

	return err
}

func GetActiveRentings(db *sql.DB) ([]int, []int, []int, []string, error) {
	query := `SELECT id FROM Renting_history WHERE end_date IS NULL`
	id, err := getMultiRowInt(db, query)

	query = `SELECT apartment_id FROM Renting_history WHERE end_date IS NULL`
	apartamentId, err := getMultiRowInt(db, query)

	query = `SELECT user_id FROM Renting_history WHERE end_date IS NULL`
	userId, err := getMultiRowInt(db, query)

	query = `SELECT start_date FROM Renting_history WHERE end_date IS NULL`
	startDate, err := getMultiRow(db, query)

	return id, apartamentId, userId, startDate, err
}

func AddNewRenting(db *sql.DB, apartamentId, userId int, startDate string) error {
	query := "INSERT INTO Renting_history (apartment_id, user_id, start_date) VALUES(?, ?, ?)"
	date, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, apartamentId, userId, date)

	return err
}

func SetEndDate(db *sql.DB, id int, endDate string) error {
	query := `UPDATE Renting_History SET end_date = ? WHERE id = ?`

	date, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, date, id)

	return err
}

func GetFaultReports(db *sql.DB) ([]int, []string, []string, []int, []int, error) {
	query := `SELECT id FROM FaultReport`
	id, err := getMultiRowInt(db, query)

	query = `SELECT description FROM FaultReport`
	description, err := getMultiRow(db, query)

	query = `SELECT date_reported FROM FaultReport`
	dateReported, err := getMultiRow(db, query)

	query = `SELECT status_id FROM FaultReport`
	statudId, err := getMultiRowInt(db, query)

	query = `SELECT apartment_id FROM FaultReport`
	apartamentId, err := getMultiRowInt(db, query)

	return id, description, dateReported, statudId, apartamentId, err
}

func AddFault(db *sql.DB, description, dateReported string, statusId, apartamentId int) error {
	query := "INSERT INTO FaultReport (description, date_reported, status_id, apartment_id) VALUES(?, ?, ?, ?)"
	date, err := time.Parse("2006-01-02", dateReported)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, description, date, statusId, apartamentId)

	return err
}

//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////

func GetEmails(db *sql.DB) ([]string, error) {
	query := `SELECT User.email FROM User`

	data, err := getMultiRow(db, query)
	return data, err

}

func GetTenents(db *sql.DB) ([]int, []string, []string, []string, []int, error) {
	query := `SELECT id FROM User WHERE role_id = "2"`
	id, err := getMultiRowInt(db, query)

	query = `SELECT name FROM User WHERE role_id = "2"`
	name, err := getMultiRow(db, query)

	query = `SELECT email FROM User WHERE role_id = "2"`
	email, err := getMultiRow(db, query)

	query = `SELECT phone FROM User WHERE role_id = "2"`
	phone, err := getMultiRow(db, query)

	query = `SELECT role_id FROM User WHERE role_id = "2"`
	role_id, err := getMultiRowInt(db, query)

	return id, name, email, phone, role_id, err
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
