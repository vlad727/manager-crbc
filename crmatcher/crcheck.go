package crcheck

import (
	"log"
	"time"
	"webapp/crmatcher/getcrname"
	"webapp/crmatcher/matchbyname"
	"webapp/crmatcher/readfile/readyamlfile"
)

var (
	SlAllowed []string
)

func CrCheck() {

	// execution time
	start := time.Now()

	log.Println("Func CrCheck started")

	// check that cluster role in yaml already exist into our allowed cr names
	if matchbyname.MatchByName(getcrname.SlAllowed, readyamlfile.Cr.Metadata.Name) {
		log.Printf("You don't need to apply thit Cluster Role yaml, because %s already exit", readyamlfile.Cr.Metadata.Name)
	} else {
		log.Printf("Cluster Role with name %s does not exist, maybe we can find something similar", readyamlfile.Cr.Metadata.Name)
	}

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for this application is %s", duration)

}
