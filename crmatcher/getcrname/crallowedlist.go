// Package getcrname get list all cluster role names
// compare it with allowed list from configmap
package getcrname

import (
	"k8s.io/utils/strings/slices"
	"log"
	"time"
	"webapp/crmatcher/getlen"
	"webapp/crmatcher/readfile/readcrifile"
)

var (
	SlAllowed []string
	MapCR     map[string]int
)

func CrAllowedList() {
	// execution time
	start := time.Now()

	// logging
	log.Println("Func CrAllowedList started")

	listAllCrs := GetCrNameList() // get all cluster roles from kubernetes API with client-go

	listForbiddenCrs := readcrifile.ReadFileCrNames() // read file with allowed cluster roles
	// iterate over slice with all cluster role names
	for _, x := range listAllCrs {
		if !slices.Contains(listForbiddenCrs, x) { // if x not in slice forbiddenCr add it to allowed slice
			SlAllowed = append(SlAllowed, x) // we will get all allowed cluster roles
		}

	}
	// map contain cluster role name and len for it Items
	MapCR = getlen.GetLen(SlAllowed)

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for func CrAllowedList is %s", duration)

}
