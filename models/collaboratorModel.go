package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Collaborator struct {
	CollaboratorID int
	UserID         int
	Username       string
	Status         string
}





func GetCollaboratorsForStory(db *sql.DB, storyID int) ([]Collaborator, error) {
	query := `
        SELECT collaborators.collaborator_id, users.id, users.username, collaborators.status
        FROM collaborators
        JOIN users ON collaborators.user_id = users.id
        WHERE collaborators.story_id = ?;
    `

	rows, err := db.Query(query, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collaborators []Collaborator
	for rows.Next() {
		var collaborator Collaborator
		err := rows.Scan(&collaborator.CollaboratorID, &collaborator.UserID, &collaborator.Username, &collaborator.Status)
		if err != nil {
			return nil, err
		}
		collaborators = append(collaborators, collaborator)
	}

	return collaborators, nil
}

func AddCollaborators(db *sql.DB, storyID int, collaboratorID int, status string) error {
	// Insert the new story into the database.
	_, err := db.Exec("INSERT INTO collaborators (story_id, user_id, status) VALUES (?, ?, ?)", storyID, collaboratorID, status)
	if err != nil {
		return err
	}
	return nil
}
