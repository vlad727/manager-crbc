package getcrb

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"net/http"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"text/template"
	"webapp/globalvar"
	"webapp/home/loggeduser"
)

// GetCrb execute after press button "Get Cluster Role Binding"
func GetCrb(w http.ResponseWriter, r *http.Request) {

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string               // name of logged user
	for k, _ := range UserAndGroups { // get logged user name from map
		username = k
	}

	// ### Read file with cluster role bindings which should hide ###
	data, err := os.ReadFile("/files/clusterroles")
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}

	dataString := string(data) // convert bytes to string

	var slCrNotAllowed []string // clear slice

	slCrNotAllowed = strings.Split(dataString, "\n") // split string and put it to slice

	//  ### Collect data to slice and map ###
	log.Println("Func GetCrb started")
	// list cluster role binding
	listCRB, err := globalvar.Clientset.RbacV1().ClusterRoleBindings().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Cluster is unavailable %s", err)
	}
	// slice for appending
	sl1 := []string{}

	// map for appending slice of strings with names
	mapTemp := map[string][]string{}
	// iterate over items to get name for cluster role binding and linked cluster role
	for _, el := range listCRB.Items {
		if slices.Contains(slCrNotAllowed, el.RoleRef.Name) {
			//log.Println("Not allowed to show ")
		} else {
			sl1 = append(sl1, "<b>"+el.Name+"</b>"+" "+el.RoleRef.Name)
			mapTemp["List"] = sl1
		}
	}
	// logging
	log.Println("Iteration over cluster role bindings finished")

	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(mapTemp)
	if err != nil {
		panic(err)
	}

	// convert to string for struct if you do not convert it will be in bytes
	str := string(yamlFile)

	//parse html
	t, _ := template.ParseFiles("tmpl/getcrb.html")

	// init struct
	Msg := struct {
		Message           string `yaml:"message"`
		MessageLoggedUser string
	}{
		Message:           str,
		MessageLoggedUser: username,
	}
	// send string to web page execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}
	// set slice to nil to prevent add new items after page refresh
	sl1 = nil

}
