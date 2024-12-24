// Package jwtdecode moved to logged user package
package jwtdecode

import (
	"github.com/golang-jwt/jwt"
	"log"
)

type MyCustomClaims struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	jwt.StandardClaims
}

// JwtDecode not used moved to logged user package
func JwtDecode(tokenData string) map[string][]string {

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
	return UserMap
}

//https://stackoverflow.com/questions/73146348/how-to-iterate-over-the-decoded-claims-of-a-jwt-token-in-go
