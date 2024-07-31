package errormsg

import (
	"net/http"
	"text/template"
	"webapp/parsepost"
)

func ErrorOut(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/error.html")
	// init struct
	Msg := struct {
		Message string
	}{
		Message: parsepost.ErrorMsg,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}

}
