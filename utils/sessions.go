package utils

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var secretKey = os.Getenv("SECRETKEY")
var store = sessions.NewCookieStore([]byte(secretKey))

func CreateSession(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := store.Get(r, "user-session")
	session.Values["username"] = username
	session.Save(r, w)
}

func GetSessionUsername(r *http.Request) (string, error) {
	session, _ := store.Get(r, "user-session")
	if username, ok := session.Values["username"].(string); ok {
		return username, nil
	}
	return "", errors.New("Session not found")
}

func LogoutSession(w http.ResponseWriter, r *http.Request) (err error) {
	session, _ := store.Get(r, "user-session")
	delete(session.Values, "username")

	session.Options = &sessions.Options{
		MaxAge:   -1, // Expire the session immediately
		HttpOnly: true,
	}
	session.Save(r, w)

	return err
}

func IsAuthenticate(r *http.Request) bool {
	user, err := GetSessionUsername(r)
	if err != nil {
		log.Println(err)
	}
	if user == "" {
		return false
	} else {
		return true
	}
}
