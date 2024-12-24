package crbcmain

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"os"
	"text/template"
	"webapp/counter"
	"webapp/getsacollect"
	"webapp/globalvar"
	"webapp/home/loggeduser"
)

type DataStruct struct {
	CrbSlice          []string
	SaMap             []string
	MessageLoggedUser string
	NumberOfEntities  int
}

func CrbcMain(w http.ResponseWriter, r *http.Request) {

	// logging
	log.Println("Func CrbcMain started")

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string               // name of logged user
	for k, _ := range UserAndGroups { // get logged user name from map
		username = k
	}
	// func Counter collect and count all cluster role bindings which one has substring "crbc"
	NumberOfEntities := counter.Counter()

	// Get Service Account and Namespaces for Select Element
	var sliceSaName []string // slice for sa name
	// why we use slice? Because we need to send it to html and add ":"
	M1, Sl1 := getsacollect.GetSaCollect(UserAndGroups)
	log.Println(M1) // got it from getsacollect func not used
	for _, x := range Sl1 {
		for k, v := range x {
			s := fmt.Sprint(k + ":" + " " + v)
			sliceSaName = append(sliceSaName, s)
		}
	}

	// Cluster Roles for Select Element
	data, err := os.ReadFile("/files/allowedlabel") // read file with allowed label
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}
	// convert bytes to string
	dataString := string(data)

	// list cluster role binding only with allowed label
	listCR, err := globalvar.Clientset.RbacV1().ClusterRoles().List(context.Background(), v1.ListOptions{LabelSelector: dataString})
	if err != nil {
		log.Println(err)
	}
	var sliceCrAllowed []string // temporary slice

	// iterate over items to get name for cluster role binding
	for _, el := range listCR.Items {
		sliceCrAllowed = append(sliceCrAllowed, el.Name)
	}
	// logging cluster roles
	log.Println("Slice cluster roles requested and collected")

	// client-go to struct for output
	DataProvider := DataStruct{
		CrbSlice:          sliceCrAllowed,   // output slice with cluster roles with allowed label
		SaMap:             sliceSaName,      // output map example: my-sa:my-namespace
		MessageLoggedUser: username,         // logged user string
		NumberOfEntities:  NumberOfEntities, // number of cluster role bindings created via manager-crbc
	}

	// parse template
	t, _ := template.ParseFiles("tmpl/crbcmain.html")

	err = t.Execute(w, DataProvider)
	if err != nil {
		return
	}

	// set slice to nil to prevent overload
	sliceSaName = nil
	sliceCrAllowed = nil
	username = ""

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	log.Println("End of crbcmain func...")

}
