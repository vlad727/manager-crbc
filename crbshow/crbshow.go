package crbshow

import (
	"net/http"
	"text/template"
	"webapp/jwtdecode"
	"webapp/parsepost"
)

func CrbShow(w http.ResponseWriter, r *http.Request) {
	//parse html
	t, _ := template.ParseFiles("tmpl/createcrbshowcrb.html")

	// init struct
	Msg := struct {
		Message           string
		MessageLoggedUser string
	}{
		Message:           parsepost.Crbname, //show created cluster role binding
		MessageLoggedUser: jwtdecode.LoggedUser,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}
}
