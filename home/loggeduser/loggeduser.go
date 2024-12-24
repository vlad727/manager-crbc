package loggeduser

import (
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

type MyCustomStruct struct {
	AuthData map[string][]string
}
type MyCustomClaims struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	jwt.StandardClaims
}

// LoggedUserRun parse request and return map with user and them groups
func LoggedUserRun(r *http.Request) map[string][]string {

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
	// get
	// authHeader := r.Header.Get("Authorization")
	// Send to Trimmer function
	JwtString := Trimmer(Data.AuthData)

	// Send to JwtDecode function
	LoggedUser := JwtDecode(JwtString)
	if len(LoggedUser) == 0 {
		log.Println("Var is empty...")
		log.Fatal("Var is empty...")
	}
	log.Printf("LoggedUser: %s", LoggedUser)
	log.Println("Func LoggedUserRun end")
	return LoggedUser
}

// Func Trimmer trim unused strings from request
func Trimmer(x map[string][]string) string {
	// temp var
	var tmp string
	// itearte over slice
	for k, v := range x {
		log.Printf("Key: %s Value: %s", k, v)
		// slice to string
		strtoken := strings.Join(v, " ")
		// delete unused part of token
		tmp = strings.ReplaceAll(strtoken, "Bearer ", "")
	}
	return tmp
}

// Func JwtDecode decode JWT token
func JwtDecode(tokenData string) map[string][]string {

	// logging
	log.Println("Func JwtDecode run")
	//log.Printf("Func JwtDecode got: %s", tokenData)
	claims := MyCustomClaims{}
	_, err := jwt.ParseWithClaims(tokenData, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("sharedKey"), nil
	})
	if err != nil {
		log.Println(err)
	}
	// logging username and groups from token
	log.Printf("LDAP username: %s", claims.Name)

	log.Printf("Groups for user: %s", claims.Groups)
	UserMap := map[string][]string{}
	// put user credentials to map
	UserMap = map[string][]string{
		claims.Name: claims.Groups,
	}
	log.Println("Func JwtDecode end")
	log.Printf("UserMap: %v", UserMap)
	return UserMap
}
