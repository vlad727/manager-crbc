package collector_bac

import (
	"fmt"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
	"webapp/globalvar"
)

var (
	// AllowedNsSlice slice for allowed namespaces
	AllowedNsSlice = []string{}
)

func Collector(m1 map[string][]string) {
	var UserName string
	var Groups []string
	// data from jwt decode
	log.Println("Got it from JWT decode: %s", m1)
	log.Println("Func Collector ....")
	for k, v := range m1 {
		UserName = k
		Groups = v
	}

	// convert slice to string
	GroupsString := fmt.Sprint(Groups)
	// get list role-bindings in namespaces
	listRB, err := globalvar.Clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed %s", listRB)
		log.Println(err)
	}
	//  logging
	log.Println(UserName)
	// iterate over role-bindings
	for _, el := range listRB.Items {
		s1 := fmt.Sprint(el.Subjects)
		// if rolebinding contain admin for roleRef and group or user contain subjects
		if el.RoleRef.Name == "admin" && strings.Contains(s1, GroupsString) || strings.Contains(s1, UserName) {

			AllowedNsSlice = append(AllowedNsSlice, el.Namespace)
		}

	}

}
