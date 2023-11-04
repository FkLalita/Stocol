package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/FkLalita/Stocol/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		// Render the registration success template
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

	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Incomplete Data", http.StatusBadRequest)
			return
		}

		// Register the user using the model
		err := models.RegisterUser(db, username, password)
		if err != nil {
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}
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
		// Replace this with your session management code

		http.Redirect(w,r, "/profile", http.StatusSeeOther)
	}       
}
    
// LogoutHandler handles user logout.                                                                                                                        
func LogoutHandler(w http.ResponseWriter, r *http.Request) {                       
	// Clear the session or cookie to log the user out
	// Replace this with your session management code   

	// Redirect to the login page or any other desired page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
