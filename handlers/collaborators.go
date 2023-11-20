package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/FkLalita/Stocol/models"
	"github.com/FkLalita/Stocol/utils"

	"github.com/gorilla/mux"
)

type CollaboratorWithStory struct {
	Collaborator models.Collaborator
	Story        models.Story
}

func RequestCollab(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	storyIDStr := vars["id"]
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		log.Println("Error converting story_id to int:", err)
		http.Error(w, "Invalid story_id parameter", http.StatusBadRequest)
		return
	}

	var user_id int
	status := "pending"
	username, err := utils.GetSessionUsername(r)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&user_id)
	if err != nil {
		log.Println("error getting user id", err)
	}
	err = models.AddCollaborators(db, storyID, user_id, status)
	if err != nil {
		log.Println("error adding collab", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func AcceptCollab(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	collabIDStr := vars["collaboratorID"]
	collaboratorID, err := strconv.Atoi(collabIDStr)
	if err != nil {
		log.Println("Error converting collaborator_id to int:", err)
		http.Error(w, "Invalid collaborator_id parameter", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		status := "Accepted"
		_, err = db.Exec("UPDATE collaborators SET status = ? WHERE collaborator_id = ?", status, collaboratorID)
		if err != nil {
			log.Println("Error updating status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func DeclineCollab(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	collabIDStr := vars["collaboratorID"]
	collaboratorID, err := strconv.Atoi(collabIDStr)
	if err != nil {
		log.Println("Error converting collaborator_id to int:", err)
		http.Error(w, "Invalid collaborator_id parameter", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		status := "Declined"
		_, err = db.Exec("UPDATE collaborators SET status = ? WHERE collaborator_id = ?", status, collaboratorID)
		if err != nil {
			log.Println("Error updating status:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetAllCollab(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var IsAuthenticated bool
	var CanEdit bool
	vars := mux.Vars(r)
	storyIDStr := vars["id"]
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		log.Println("Error converting story_id to int:", err)
		http.Error(w, "Invalid story_id parameter", http.StatusBadRequest)
		return
	}
	user, err := utils.GetSessionUsername(r)
	if user == "" {
		IsAuthenticated = false
	} else {
		IsAuthenticated = true
	}

	
	if Collaborator. == "" {
		IsAuthenticated = false
	} else {
		IsAuthenticated = true
	}
	
	Collaborators, err := models.GetCollaboratorsForStory(db, storyID)
	if err != nil {
		log.Println("Error getting collabs:", err)
	}

	Story, err := models.GetStoryDetails(db, storyID)
	if err != nil {
		log.Println("Error fetching story details:", err)
	}



	var collaboratorsWithStory []CollaboratorWithStory
	for _, collaborator := range Collaborators {
		collaboratorWithStory := CollaboratorWithStory{
			Collaborator: collaborator,
			Story:        Story,
		}
		collaboratorsWithStory = append(collaboratorsWithStory, collaboratorWithStory)
	}

	data := struct {
		IsAuthenticated bool
		Collaborators   []CollaboratorWithStory
	}{
		IsAuthenticated: IsAuthenticated,
		Collaborators:   collaboratorsWithStory,
	}

	tmpl, err := template.ParseFiles("templates/collaborator.html", "templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("error executing temp:", err)
	}
}
