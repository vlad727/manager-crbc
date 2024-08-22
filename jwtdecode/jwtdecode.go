package jwtdecode

import (
	"github.com/golang-jwt/jwt"
	"log"
)

var (
	UserMap map[string][]string

	// LoggedUser var for web ui to show which user logged in
	LoggedUser string
)

type MyCustomClaims struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	jwt.StandardClaims
}

func JwtDecode(tokenData string) {

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
	// var for web ui to show which user logged in
	LoggedUser = claims.Name
	log.Printf("Groups for user: %s", claims.Groups)
	// put user credentials to map
	UserMap = map[string][]string{
		claims.Name: claims.Groups,
	}

}

//https://stackoverflow.com/questions/73146348/how-to-iterate-over-the-decoded-claims-of-a-jwt-token-in-go
