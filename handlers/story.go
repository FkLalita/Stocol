package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/FkLalita/Stocol/models"
	"github.com/FkLalita/Stocol/utils"

	"github.com/gorilla/mux"
)

func CreateStoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		var IsAuthenticated bool
		IsAuthenticated = utils.IsAuthenticate(r)
		data := struct {
			IsAuthenticated bool
		}{
			IsAuthenticated: IsAuthenticated,
		}

		tmpl, err := template.ParseFiles("templates/create_story.html", "templates/base.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodPost {
		var authorID int
		var authorName string

		username, err := utils.GetSessionUsername(r)
		if err != nil {
			log.Println("Error getting session username:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		created_at := time.Now()

		// Query for author details
		err = db.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&authorID, &authorName)
		if err != nil {
			log.Println("Error fetching author details:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if title == "" || content == "" {
			http.Error(w, "Incomplete Data", http.StatusBadRequest)
			return
		}

		// Create the story with the model
		if err := models.CreateStory(db, title, content, created_at, authorID, authorName); err != nil {
			log.Println("Error creating story:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func ViewStory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var IsAuthenticated bool
	var IsAuthor bool
	vars := mux.Vars(r)
	storyIDStr := vars["id"]
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		log.Println("Error converting story_id to int:", err)
		http.Error(w, "Invalid story_id parameter", http.StatusBadRequest)
		return
	}

	Story, err := models.GetStoryDetails(db, storyID)
	if err != nil {
		log.Println("Error fetching story details:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	username, err := utils.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "Error Fetching Username", http.StatusBadRequest)
	}

	if username == Story.AuthorName {
		IsAuthor = true
	} else {
		IsAuthor = false
	}

	IsAuthenticated = utils.IsAuthenticate(r)
	data := struct {
		IsAuthor        bool
		Story           models.Story
		IsAuthenticated bool
	}{

		IsAuthor:        IsAuthor,
		Story:           Story,
		IsAuthenticated: IsAuthenticated,
	}

	tmpl, err := template.ParseFiles("templates/view_story.html", "templates/base.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with the story data

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("Error executing template:", err)

	} 

}
func UpdateStory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	storyIDStr := vars["id"]
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		log.Println("Error converting story_id to int:", err)
		http.Error(w, "Invalid story_id parameter", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		// Retrieve form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		created_at := time.Now()

		// Get author details
		username, err := utils.GetSessionUsername(r)
		if err != nil {
			log.Println("Error getting session username:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var authorID int
		var authorName string
		err = db.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&authorID, &authorName)
		if err != nil {
			log.Println("Error fetching author details:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Update the story with the model
		err = models.UpdateStory(db, storyID, title, content, created_at, authorID, authorName)
		if err != nil {
			log.Println("Error updating story:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method != http.MethodPost {
		var IsAuthenticated bool
		var IsAuthor bool

		Story, err := models.GetStoryDetails(db, storyID)
		if err != nil {
			log.Println("Error fetching story details:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		username, err := utils.GetSessionUsername(r)
		if err != nil {
			http.Error(w, "Error Fetching Username", http.StatusBadRequest)
			return
		}

		if username == Story.AuthorName {
			IsAuthor = true
		} else {
			IsAuthor = false
		}

		IsAuthenticated = utils.IsAuthenticate(r)

		data := struct {
			IsAuthor        bool
			Story           models.Story
			IsAuthenticated bool
		}{
			IsAuthor:        IsAuthor,
			Story:           Story,
			IsAuthenticated: IsAuthenticated,
		}

		tmpl, err := template.ParseFiles("templates/update_story.html", "templates/base.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template with the story data
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Println("Error executing template:", err)
		}
	}
}

func DeleteStoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	storyIDStr := vars["id"]
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		log.Println("Error converting story_id to int:", err)
		http.Error(w, "Invalid story_id parameter", http.StatusBadRequest)
		return
	}

	// Call a function to delete the story by ID
	err = models.DeleteStory(db, storyID)
	if err != nil {
		log.Println("Error deleting story:", err)
		http.Error(w, "Error Deleting Story", http.StatusInternalServerError)
		return
	}

	// Redirect to the home page or another appropriate page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
