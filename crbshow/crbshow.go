package crbshow

import (
	"net/http"
	"text/template"
	"webapp/parsepost"
)

func CrbShow(w http.ResponseWriter, r *http.Request) {
	//parse html
	t, _ := template.ParseFiles("tmpl/createcrbshowcrb.html")

	// init struct
	Msg := struct {
		Message string
	}{
		Message: parsepost.Crbname, //show created cluster role binding
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}
}
