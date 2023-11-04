package handlers

import (
	"html/template"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)

}



