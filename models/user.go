package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   int
	Username string
	Password string
}

// RegisterUser adds a new user to the database.
func RegisterUser(db *sql.DB, username, password string) error {
	// Hash the password before storing it in the database.
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	// Insert the new user into the database.
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return err

	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyUser verifies a user's credentials.
func VerifyUser(db *sql.DB, username, password string) error {
	// Retrieve the hashed password from the database for the given username.
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	// Compare the provided password with the stored hashed password.
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
