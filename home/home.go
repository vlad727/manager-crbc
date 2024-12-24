package home

import (
	"net/http"
	"text/template"
	"webapp/home/loggeduser"
)

// HomeFunc the main page
func HomeFunc(w http.ResponseWriter, r *http.Request) {

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string               // empty var for name of logged user
	for k, _ := range UserAndGroups { // get logged user name from map and skip groups
		username = k
	}
	// parse html
	t, _ := template.ParseFiles("tmpl/getresp.html")

	// create and init struct
	UserStruct := struct {
		Message string
	}{
		Message: username, // set logged user login and put to html
	}

	err := t.Execute(w, UserStruct)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	// set string with name to nil
	username = ""

}
