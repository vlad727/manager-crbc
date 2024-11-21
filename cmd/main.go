package main

import (
	"log"
	"net/http"
	"time"
	"webapp/crbcmain"
	"webapp/crbshow"
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
	//	porthttps = ":8443" webhook won't use with this application
)

var (
	tlscert, tlskey string
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
	http.HandleFunc("/crbshow", crbshow.CrbShow)                     // show result after creating cluster role binding
	http.HandleFunc("/error", errormsg.ErrorOut)                     // show page with error
	http.HandleFunc("/getcrdesc", getcrdesc.GetCrDesc)               // it get post request parse and redirect to page with result
	http.HandleFunc("/uploadfile", handlers.UploadFile)              // crmatcher -> handlers ->  handlerpost.go upload file
	http.HandleFunc("/uploadedfile", handlers.HandlePost)            // already uploaded file send to parse data from file
	http.HandleFunc("/crmatcherresult", handlers.CrMatcherResult)    // crmatcherresult -> handlers ->  handlerresult.go show page with result checking cluster roles
	//http.HandleFunc("/health", checkhealth.Health)                   // allow check health for application
	/* webhook won't use with this application
	// goroutine for webhook part port 8443
	// need to run it in goroutine because http.ListenAndServe can't listen on two ports at the same time
	go func() {

		log.Printf("Port %s listening", porthttps)
		flag.StringVar(&tlscert, "tlsCertFile", "/certs/tls.crt",
			"File containing a certificate for HTTPS.")
		flag.StringVar(&tlskey, "tlsKeyFile", "/certs/tls.key",
			"File containing a private key for HTTPS.")
		flag.Parse()
		// func validate.bac in package webhook
		http.HandleFunc("/validate.bac", validate.bac.Validate)
		log.Fatal(http.ListenAndServeTLS(porthttps, tlscert, tlskey, nil))
	}()

	*/

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func main  %s", duration)

	// listen http
	http.ListenAndServe(porthttp, nil)

}
