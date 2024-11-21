package getsa

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"net/http"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/globalvar"
	"webapp/jwtdecode"
	"webapp/readfiles"
)

var (
	// AllowedNsSlice slice for allowed namespaces
	AllowedNsSlice = []string{}

	// temp string for name from RB Subject
	strNameFromSub string

	// var for data from jwtdecode
	UserName  string
	Groups    []string
	UserAdmin string
)

type StructGetSa struct {
	Collection map[string]string
}

func GetSa(w http.ResponseWriter, r *http.Request) {

	// Run func for read file with user admin
	UserAdmin = readfiles.ReadFile()

	// get len for var and string
	log.Println("Get Len")
	log.Println(len(UserAdmin))
	log.Println(len("admin"))

	log.Printf("Got it from func read file %s", UserAdmin)

	// logging
	log.Println("Func GetSa ....")
	// data from jwt decode
	//log.Println("Got it from JWT decode: %s", jwtdecode.UserMap)

	// iterate over map to assign data to new clientgo.bac
	for k, v := range jwtdecode.UserMap { // clientgo.bac comes from jwtdecode func
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
		// readfiles.UserAdmin <-- get var from configmap func ReadFile
		if el.RoleRef.Name == UserAdmin && strNameFromSub == UserName || slices.Contains(Groups, strNameFromSub) {
			AllowedNsSlice = append(AllowedNsSlice, el.Namespace)
		}

	}
	// logging to know which one namespace we got
	log.Printf("Allowed namespaces: %s", AllowedNsSlice)

	// temporary slice for service accounts
	tmpSl := []string{}
	// main slice for output
	slNsSa := []map[string][]string{}
	// map for key=ns-name value=[sa names]
	m1 := make(map[string][]string)
	// iterate over service accounts
	for _, y := range AllowedNsSlice { // slice namespaces from collector.bac
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
		Message           string `yaml:"message"`
		MessageLoggedUser string
	}{
		Message:           str,
		MessageLoggedUser: jwtdecode.LoggedUser,
	}
	// execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}

	// set slice to nil to avoid repeat namespaces after refresh page
	AllowedNsSlice = nil
}
