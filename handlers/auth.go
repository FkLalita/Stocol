package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/FkLalita/Stocol/models"
	"github.com/FkLalita/Stocol/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		// Render the registration form
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")

		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Incomplete Data", http.StatusBadRequest)
			return
		}

		// Check if the username already exists
		var existingUsername string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
		if err == nil {
			// Username already exists
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		} else if err != sql.ErrNoRows {
			// An error occurred while checking the database
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Register the user using the model
		err = models.RegisterUser(db, username, password)
		if err != nil {
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}

		// Set a session or cookie to indicate that the user is logged in
		utils.CreateSession(w, r, username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Incomplete Data", http.StatusBadRequest)
			return
		}

		// Verify user credentials using the model
		if err := models.VerifyUser(db, username, password); err != nil {
			http.Error(w, "Login failed", http.StatusUnauthorized)
			log.Println(err)
			return
		}

		// Set a session or cookie to indicate that the user is logged in
		utils.CreateSession(w, r, username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// LogoutHandler handles user logout.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := utils.LogoutSession(w, r)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
