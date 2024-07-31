package crbcmain

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"text/template"
	"webapp/globalvar"
)

type DataStruct struct {
	TestData []string
	CrbSlice []string
	SaMap    map[string]string
}

func CrbcMain(w http.ResponseWriter, r *http.Request) {
	log.Println("Func CrbcMain started")
	//---------------------------------------------------------------------------------------------------------------------------------
	// Service Account and their namespaces
	m := map[string]string{}
	log.Println("Collect service accounts items")
	listSa, err := globalvar.Clientset.CoreV1().ServiceAccounts("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Cluster is unavailable %s", err)

	} else {
		log.Println("Requested service account list from API")
	}
	// iterate over items to get sa and ns and put to string use for it string builder
	for _, el := range listSa.Items {

		//log.Println(el.Namespace, el.Name)

		m[el.Name] = el.Namespace

	}
	// logging service accounts
	//log.Println(m)
	log.Println("Map with service accounts and their namespaces has been created")
	//---------------------------------------------------------------------------------------------------------------------------------
	// Cluster Roles part
	// init map
	sl := []string{}
	// list cluster role binding
	listCR, err := globalvar.Clientset.RbacV1().ClusterRoles().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	// iterate over items to get name for cluster role binding and linked cluster role
	for _, el := range listCR.Items {
		// Logging Cluster Roles
		//log.Printf(" ClusterRole: %s ", el.Name)
		sl = append(sl, el.Name)

	}
	// logging cluster roles
	//log.Printf("Console output for cluster role: %s", sl)
	log.Println("Slice cluster roles requested and collected")

	//---------------------------------------------------------------------------------------------------------------------------------

	DataProvider := DataStruct{

		TestData: []string{ // for test
			"Extra Priority",
			"Normal",
			"Low Priority"},
		CrbSlice: sl, // output slice
		SaMap:    m,  // output map
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	t, _ := template.ParseFiles("tmpl/crbcmain.html")

	err = t.Execute(w, DataProvider)
	if err != nil {
		return
	}

}
