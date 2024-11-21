// Package readfiles read data from file forbiddencrs convert to string and return slice to main
package readcrifile

import (
	"log"
	"os"
	"strings"
	"time"
)

var (
	sl []string
)

func ReadFileCrNames() []string {

	start := time.Now()
	// logging readFile
	log.Println("Func ReadFileCrNames started")

	// read file with user admin
	crNames, err := os.ReadFile("/files/clusterroles")
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}
	strCrNames := string(crNames)
	s := strings.Split(strCrNames, "\n")
	for _, x := range s {
		sl = append(sl, x)
	}
	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func ReadFileCrNames  %s", duration)

	return sl
}
