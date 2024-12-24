package errormsg

import (
	"log"
	"net/http"
	"text/template"
	"webapp/home/loggeduser"
)

func ErrorOut(w http.ResponseWriter, r *http.Request) {

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string               // name of logged user
	for k, _ := range UserAndGroups { // get logged user name from map
		username = k
	}

	errorMessage := r.URL.Query().Get("error") // get value for key "error"
	log.Printf("Message from errorMessage var %s", errorMessage)
	t, _ := template.ParseFiles("tmpl/error.html")
	// init struct
	Msg := struct {
		Message           string
		MessageLoggedUser string
	}{
		Message:           errorMessage,
		MessageLoggedUser: username,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}

}
