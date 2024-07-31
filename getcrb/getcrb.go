package getcrb

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/globalvar"
)

func GetCrb(w http.ResponseWriter, r *http.Request) {
	// ---------------------------------------------------------------------------------------------------------
	// collect data to slice and map
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

		sl1 = append(sl1, "<b>"+el.Name+"</b>"+" "+el.RoleRef.Name)
		mapTemp["List"] = sl1

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
		Message string `yaml:"message"`
	}{
		Message: str,
	}
	// send string to web page execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}
	// set slice to nil to prevent add new items after page refresh
	sl1 = nil

}
