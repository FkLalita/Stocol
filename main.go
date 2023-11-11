package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/FkLalita/Stocol/handlers"
	"github.com/FkLalita/Stocol/models"
	"github.com/FkLalita/Stocol/utils"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "FkLalita:ayomide.10@tcp(localhost:3306)/Stocol")
	if err != nil {
		fmt.Println("Error Conneting To Database", err)
		return
	}

	defer db.Close()

	// Define routes and handlers for user registration, authentication, and story collaboration.
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(db, w, r)
	})

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, db)
	})
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db)
	})

	r.HandleFunc("/profile", handlers.ProfileHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler)

	r.HandleFunc("/create-story", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateStoryHandler(w, r, db)
	})
	r.HandleFunc("/story/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handlers.ViewStory(w, r, db)
	})

	r.HandleFunc("/story/{id:[0-9]+}/edit", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateStory(w, r, db)
	})

	r.HandleFunc("/story/{id:[0-9]+}/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteStoryHandler(w, r, db)
	})

	// Start the server.
	fmt.Println("Starting Server.................")
	http.ListenAndServe(":8080", r)
}

func homeHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var IsAuthenticated bool
	user, err := utils.GetSessionUsername(r)
	if user == "" {
		IsAuthenticated = false
	} else {
		IsAuthenticated = true
	}

	stories := models.GetAllStories(db)
	data := struct {
		IsAuthenticated bool
		Stories         []models.Story
	}{
		IsAuthenticated: IsAuthenticated,
		Stories:         stories,
	}

	tmpl, err := template.ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)

}
