package readfiles

import (
	"log"
	"os"
	"strings"
)

func ReadFile() string {

	// logging readFile
	log.Println("Func ReadFile started")
	// read file with user admin
	data, err := os.ReadFile("/files/user-admin")
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}
	// declare variable
	var UserAdmin = ""
	// convert bytes to string
	UserAdmin = string(data)
	// logging
	log.Printf("Got username %s", UserAdmin)
	tmp := strings.ReplaceAll(UserAdmin, "\n", "")
	// return string
	return tmp

}
