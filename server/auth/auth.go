package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"server/database"
)

func GenerateSecureToken(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

func CreateSession(db *sql.DB, userId int) (string, error) {
	token, err := GenerateSecureToken(64)
	if err != nil {
		//log.Println(err)
		return "", err
	}

	if err = database.InsertToken(db, userId, token); err != nil {
		//log.Println(err)
		return "", nil
	}

	//log.Println(token)
	return token, nil
}

func ValidateSession(db *sql.DB, token string) (int, error) {
	userId, err := database.GetToken(db, token)

	return userId, err
}
