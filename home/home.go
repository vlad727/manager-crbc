package home

import (
	"log"
	"net/http"
	"text/template"
	"webapp/jwtdecode"
	"webapp/trimmer"
)

var (
	LoggedUser string
)

type MyCustomStruct struct {
	AuthData map[string][]string
}

func HomeFunc(w http.ResponseWriter, r *http.Request) {

	var Data MyCustomStruct

	r.ParseForm() // Анализирует переданные параметры url, затем анализирует пакет ответа для тела POST (тела запроса)
	// внимание: без вызова метода ParseForm последующие данные не будут получены
	//log.Println(r.Header)
	// Loop over header names
	for name, values := range r.Header {

		//log.Println(name, values)
		if name == "Authorization" {
			//log.Printf("Token with string \"Bearer\" need to trim %s", values)
			Data = MyCustomStruct{
				AuthData: map[string][]string{
					"Authorization": values,
				},
			}

		}

	}
	// send to extract token
	jwtToken := trimmer.Trimmer(Data.AuthData)

	// send token to decode
	jwtdecode.JwtDecode(jwtToken)
	// logging
	log.Println("Func HomeFunc started")

	// parse html
	t, _ := template.ParseFiles("tmpl/home.html")

	// create and init struct
	Marketing := struct {
		Message string
	}{
		Message: jwtdecode.LoggedUser, // get logged user name from jwt decode
	}

	err := t.Execute(w, Marketing)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

}
