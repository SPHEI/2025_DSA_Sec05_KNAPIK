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
