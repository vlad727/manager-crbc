package home

import (
	"log"
	"net/http"
	"text/template"
)

func HomeFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Func HomeFunc started")
	t, _ := template.ParseFiles("tmpl/home.html")

	Marketing := struct {
		Message string
	}{
		Message: "Cluster Role Binding Creator",
	}
	err := t.Execute(w, Marketing)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
}
