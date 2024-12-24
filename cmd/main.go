package main

import (
	"log"
	"net/http"
	"time"
	"webapp/checkhealth"
	"webapp/crbcmain"
	"webapp/crmatcher/getcrname"
	"webapp/crmatcher/handlers"
	errormsg "webapp/error"
	"webapp/getcrb"
	"webapp/getcrdesc"
	"webapp/getsa"
	"webapp/home"
	"webapp/parsepost"
)

const (
	porthttp = ":8080"
)

func main() {

	// // Code to measure
	start := time.Now()

	getcrname.CrAllowedList() // func which one count len for items all allowed cluster roles

	// logging
	log.Println("Hello my dear friend")
	log.Println("I see you like to press buttons")
	log.Printf("Port %s listening", porthttp)
	log.Println("Func main started")

	// All handlers for get/post requests
	http.HandleFunc("/", home.HomeFunc)                              // home page with buttons main page for application
	http.HandleFunc("/getcrb", getcrb.GetCrb)                        // allow getting cluster role binding as a list
	http.HandleFunc("/getsa", getsa.GetSa)                           // allow getting service accounts and their namespaces
	http.HandleFunc("/crbcmain", crbcmain.CrbcMain)                  // generate page with fields allow choosing service account ns and cluster role
	http.HandleFunc("/createcrbmanager", parsepost.ParsePostRequest) // parse input from user service account + namespace + cluster role + crbc
	//http.HandleFunc("/crbshow", crbshow.CrbShow)                     // show result after creating cluster role binding // deprecated see parsepostrequest and error packages
	http.HandleFunc("/error", errormsg.ErrorOut)                  // show page with error
	http.HandleFunc("/getcrdesc", getcrdesc.GetCrDesc)            // it get post request parse and redirect to page with result
	http.HandleFunc("/uploadfile", handlers.UploadFile)           // crmatcher -> handlers ->  handlerpost.go upload file
	http.HandleFunc("/uploadedfile", handlers.HandlePost)         // already uploaded file send to parse data from file
	http.HandleFunc("/crmatcherresult", handlers.CrMatcherResult) // crmatcherresult -> handlers ->  handlerresult.go show page with result checking cluster roles
	http.HandleFunc("/health", checkhealth.Health)                // allow check health for application

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func main  %s", duration)

	// listen http
	http.ListenAndServe(porthttp, nil)
}
