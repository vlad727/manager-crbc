package getcrdesc

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/globalvar"
	"webapp/jwtdecode"
)

var (
	ClusterRoleName = ""
)

func GetCrDesc(w http.ResponseWriter, r *http.Request) {

	// parse post request
	r.ParseForm() // Анализирует переданные параметры url, затем анализирует пакет ответа для тела POST (тела запроса)
	// внимание: без вызова метода ParseForm последующие данные не будут получены
	//log.Printf("Full post request: %s", r)
	log.Println(r.Form) // печатает информацию на сервере
	log.Println("Path: ", r.URL.Path)

	// iterate over request with for to get requested cluster role name
	for k, v := range r.Form {
		log.Println("Key: ", k)
		log.Printf("Value:%s", v)
		for _, el := range v {
			ClusterRoleName = el
		}

	}
	// logging
	log.Println("Func GetCrDesc ...")

	// get cluster role description
	getCr, err := globalvar.Clientset.RbacV1().ClusterRoles().Get(context.Background(), ClusterRoleName, v1.GetOptions{})
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Cluster is unavailable %s", err)
	}
	// main slice
	outSlice := []map[string][]string{}

	// temp slices
	sl0 := []string{}
	sl1 := []string{}
	sl2 := []string{}
	sl3 := []string{}
	sl4 := []string{}

	// temp maps
	m0 := map[string][]string{}
	m1 := map[string][]string{}
	m2 := map[string][]string{}
	m3 := map[string][]string{}
	m4 := map[string][]string{}

	// iterate over items to get sa and ns and put it to string with string builder
	for _, el := range getCr.Rules {

		tempslice := [][]string{el.APIGroups, el.ResourceNames, el.Resources, el.Verbs, el.NonResourceURLs}
		for x, item := range tempslice {

			switch x {
			case 0:
				for _, a := range item {
					if len(a) == 0 {
						a += "\"\""
					}
					sl0 = append(sl0, a)

				}
				// add data to map
				m0["apiGroups"] = sl0
				// add data to main slice
				outSlice = append(outSlice, m0)

				// clear map
				m0 = make(map[string][]string)
				// clear slice
				sl0 = nil

			case 1:
				if len(item) == 0 {
					log.Println("No resourceNames")
				} else {
					for _, a := range item {

						sl1 = append(sl1, a)

					}
					m1["resourceNames"] = sl1
					outSlice = append(outSlice, m1)

					// clear map
					m1 = make(map[string][]string)
					// clear slice
					sl1 = nil
				}
			case 2:
				for _, a := range item {

					sl2 = append(sl2, a)

				}
				m2["resources"] = sl2
				outSlice = append(outSlice, m2)

				// clear map
				m2 = make(map[string][]string)
				// clear slice
				sl2 = nil
			case 3:
				for _, a := range item {

					sl3 = append(sl3, a)

				}
				m3["verbs"] = sl3
				outSlice = append(outSlice, m3)
				// clear map
				m3 = make(map[string][]string)
				// clear slice
				sl3 = nil
			case 4:

				if len(item) == 0 {
					log.Println("No nonResourceURLs")
				} else {
					for _, a := range item {

						sl4 = append(sl4, a)

					}
					m4["nonResourceURLs"] = sl4
					outSlice = append(outSlice, m4)
					// clear map
					m4 = make(map[string][]string)
					// clear slice
					sl4 = nil
				}

			}

		}

	}
	// logging slice
	log.Println(outSlice)

	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(outSlice)
	if err != nil {
		panic(err)
	}

	// convert to string
	s := string(yamlFile)

	// init struct with var
	Msg := struct {
		ClusterRoleName   string
		Items             string `yaml:"out"`
		MessageLoggedUser string
	}{
		ClusterRoleName:   ClusterRoleName,
		Items:             s,
		MessageLoggedUser: jwtdecode.LoggedUser,
	}

	// parse html
	t, _ := template.ParseFiles("tmpl/descshow.html")

	// execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}

	// redirect to page with description
	http.Redirect(w, r, "/descshow", http.StatusSeeOther)
}
