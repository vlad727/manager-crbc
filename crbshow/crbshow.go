// Package deprecated
package crbshow

import (
	"net/http"
	"text/template"
	"webapp/home/loggeduser"
)

func CrbShow(w http.ResponseWriter, r *http.Request) {
	//parse html
	t, _ := template.ParseFiles("tmpl/createcrbshowcrb.html")
	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string               // name of logged user
	for k, _ := range UserAndGroups { // get logged user name from map
		username = k
	}

	// init struct
	Msg := struct {
		Message           string
		MessageLoggedUser string
	}{
		//	Message:           parsepost.Crbname, //show created cluster role binding
		MessageLoggedUser: username,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}
}
