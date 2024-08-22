package crbcmain

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"webapp/globalvar"
	"webapp/jwtdecode"
	"webapp/readfiles"
)

type DataStruct struct {
	TestData []string
	CrbSlice []string
	SaMap    []string
}

var (
	// AllowedNsSlice slice for allowed namespaces
	AllowedNsSlice = []string{}

	// temp string for name from RB Subject
	strNameFromSub string

	// var for data from jwtdecode
	UserName string
	Groups   []string

	sliceSaName    []string
	sliceCrAllowed []string
	UserAdmin      string
)

func CrbcMain(w http.ResponseWriter, r *http.Request) {

	// Run func for ReadFile to get value from config file
	UserAdmin = readfiles.ReadFile()

	log.Println("Func CrbcMain started")

	// data from jwt decode
	log.Println("Got it from JWT decode: %s", jwtdecode.UserMap)

	// iterate over map to assign data to new vars
	for k, v := range jwtdecode.UserMap { // vars comes from jwtdecode func
		UserName = k
		Groups = v
	}

	// get list role-bindings in namespaces
	listRB, err := globalvar.Clientset.RbacV1().RoleBindings("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("Failed %s", listRB)
		log.Println(err)
	}

	// iterate over role-bindings
	for _, el := range listRB.Items {
		// iterate over Subjects to get name (also it contains: apiGroup, kind, namespace )
		for _, x := range el.Subjects {
			//log.Println(x.Name)
			strNameFromSub = x.Name //May be group or username from ldap

		}
		// check condition: if clusterRole == admin and linked with user or group add namespace to allowed list
		if el.RoleRef.Name == UserAdmin && strNameFromSub == UserName || slices.Contains(Groups, strNameFromSub) {
			AllowedNsSlice = append(AllowedNsSlice, el.Namespace)
		}

	}
	// logging to know which one namespace we got
	//log.Printf("Allowed namespaces: %s", AllowedNsSlice)

	//---------------------------------------------------------------------------------------------------------------------------------
	// collect service accounts and their namespaces
	// ---------------------------------------------------------------------------------------------------------------------
	for _, y := range AllowedNsSlice {
		listSa, err := globalvar.Clientset.CoreV1().ServiceAccounts(y).List(context.Background(), v1.ListOptions{})
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Cluster is unavailable %s", err)

		} else {
			log.Println("Requested service account list from API")
		}

		//m := map[string][]string{}
		for _, el := range listSa.Items {
			s := fmt.Sprint(el.Namespace + " " + el.Name)
			sliceSaName = append(sliceSaName, s)
		}
	}

	//---------------------------------------------------------------------------------------------------------------------------------
	// Cluster Roles part
	//---------------------------------------------------------------------------------------------------------------------------------
	slCrNotAllowed := []string{}

	// read file with cluster roles which one should hide
	data, err := os.ReadFile("/files/clusterroles")
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}
	// convert bytes to string
	dataString := string(data)

	// split string and put it to slice
	slCrNotAllowed = strings.Split(dataString, "\n")

	// logging slice to know what we got
	//log.Println(slCrNotAllowed)

	// list cluster role binding
	listCR, err := globalvar.Clientset.RbacV1().ClusterRoles().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	// iterate over items to get name for cluster role binding and linked cluster role
	for _, el := range listCR.Items {

		// check cluster role does it present in forbidden list
		if slices.Contains(slCrNotAllowed, el.Name) {
			log.Printf("Cluster Role %s should hide", el.Name)
		} else {
			// Logging Cluster Roles
			//log.Printf(" ClusterRole: %s ", el.Name)
			sliceCrAllowed = append(sliceCrAllowed, el.Name)
		}

	}
	// logging cluster roles
	log.Println("Slice cluster roles requested and collected")

	//---------------------------------------------------------------------------------------------------------------------------------
	// vars to struct
	DataProvider := DataStruct{

		TestData: []string{ // for test
			"Extra Priority",
			"Normal",
			"Low Priority"},
		CrbSlice: sliceCrAllowed, // output slice
		SaMap:    sliceSaName,    // output map
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	t, _ := template.ParseFiles("tmpl/crbcmain.html")

	err = t.Execute(w, DataProvider)
	if err != nil {
		return
	}

	// set slice to nil
	sliceSaName = nil
	AllowedNsSlice = nil
	sliceCrAllowed = nil

}
