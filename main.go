package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/FkLalita/Stocol/handlers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "FkLalita:ayomide.10@tcp(localhost:3306)/stocol")
	if err != nil {
		fmt.Println("Error Conneting To Database", err)
		return
	}

	defer db.Close()

	// Define routes and handlers for user registration, authentication, and story collaboration.
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)

	// Start the server.
	fmt.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}





func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)

}
