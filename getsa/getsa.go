package getsa

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

type StructGetSa struct {
	Collection map[string]string
}

func GetSa(w http.ResponseWriter, r *http.Request) {

	// ---------------------------------------------------------------------------------------------------------
	// empty slice for namespaces
	slNs := []string{}
	// get namespaces names
	listNs, err := globalvar.Clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Cluster is unavailable %s", err)
	}
	// iterate over namespace name and put it to slice
	for _, x := range listNs.Items {
		slNs = append(slNs, x.Name)
	}
	// temporary slice for service accounts
	tmpSl := []string{}
	// main slice for output
	slNsSa := []map[string][]string{}
	// map for key=ns-name value=[sa names]
	m1 := make(map[string][]string)
	// iterate over service accounts
	for _, y := range slNs {
		// get service account list and their namespaces
		listSa, _ := globalvar.Clientset.CoreV1().ServiceAccounts(y).List(context.Background(), v1.ListOptions{})
		// iterate over service aacounts
		for _, z := range listSa.Items {
			tmpSl = append(tmpSl, z.Name)

		}
		// set key ns name + slice service account
		m1[y] = tmpSl
		// append map to main slice
		slNsSa = append(slNsSa, m1)
		// clear slice
		tmpSl = nil
		// clear map
		m1 = make(map[string][]string)
	}
	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(slNsSa)
	if err != nil {
		panic(err)
	}
	// convert to string for struct if you do not convert it will be in bytes
	str := string(yamlFile)
	// parse html
	t, _ := template.ParseFiles("tmpl/getsa.html")
	// init struct and var
	Msg := struct {
		Message string `yaml:"message"`
	}{
		Message: str,
	}
	// execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}
}
