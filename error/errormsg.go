package errormsg

import (
	"net/http"
	"text/template"
	"webapp/jwtdecode"
	"webapp/parsepost"
)

func ErrorOut(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/error.html")
	// init struct
	Msg := struct {
		Message           string
		MessageLoggedUser string
	}{
		Message:           parsepost.ErrorMsg,
		MessageLoggedUser: jwtdecode.LoggedUser,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}

}
