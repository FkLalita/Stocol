package handlers

import (
	"html/template"
	"net/http"

	"github.com/FkLalita/Stocol/utils"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	Username, err := utils.GetSessionUsername(r)
	
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, Username)

}
