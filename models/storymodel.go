package models

import (
	"database/sql"
	"log"
	"time"
)

type Story struct {
	StoryID    int
	Title      string
	Content    string
	CreatedAt  time.Time
	AuthorID   int
	AuthorName string
}

func CreateStory(db *sql.DB, title, content string, created_at time.Time, authorID int, authorName string) error {
	// Insert the new story into the database.
	_, err := db.Exec("INSERT INTO stories (title, content, created_at, author_id, author_name) VALUES (?, ?, ?, ?, ?)", title, content, created_at, authorID, authorName)
	if err != nil {
		return err
	}
	return nil
}

func GetAllStories(db *sql.DB) []Story {
	var stories []Story

	rows, err := db.Query("SELECT * FROM stories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var s Story
		var createdAtStr []uint8 // Temporary variable to store the string from the database.

		if err := rows.Scan(&s.StoryID, &s.Title, &s.Content, &createdAtStr, &s.AuthorID, &s.AuthorName); err != nil {
			log.Println(err)
			continue
		}

		// Convert the createdAtStr ([]uint8) to a string.
		createdAtString := string(createdAtStr)

		// Parse the createdAtString as time.Time.
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtString)
		if err != nil {
			log.Println(err)
		} else {
			s.CreatedAt = parsedTime
		}

		stories = append(stories, s)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return stories
}


func GetStoryDetails(db *sql.DB, story_id int) (Story, error) {
    var s Story

    // Query the database to get story details by story ID
    row := db.QueryRow("SELECT * FROM stories WHERE story_id = ?", story_id)

    // Temporary variable to store the string from the database.
    var createdAtStr []uint8

    // Scan the row to retrieve story details
    if err := row.Scan(&s.StoryID, &s.Title, &s.Content, &createdAtStr, &s.AuthorID, &s.AuthorName); err != nil {
        // Handle the error, log it, and return an error
        log.Println(err)
        return Story{}, err
    }

    // Convert the createdAtStr ([]uint8) to a string.
    createdAtString := string(createdAtStr)

    // Parse the createdAtString as time.Time.
    parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtString)
    if err != nil {
        // Handle the error, log it, and return an error
        log.Println(err)
        return Story{}, err
    }

    // Set the CreatedAt field of the story
    s.CreatedAt = parsedTime

    // Return the retrieved story
    return s, nil
}

func UpdateStory(db *sql.DB, storyID int, title, content string, createdAt time.Time, authorID int, authorName string) error {
    _, err := db.Exec("UPDATE stories SET title=?, content=?, created_at=?, author_id=?, author_name=? WHERE story_id=?", title, content, createdAt, authorID, authorName, storyID)
    return err
}












func DeleteStory(db *sql.DB, storyID int) error {
    _, err := db.Exec("DELETE FROM stories WHERE story_id = ?", storyID)
    return err
}

