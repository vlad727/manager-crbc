// Package getlen iterate over all cluster role and their items to get len for all rules
// example and the end of code
package getlen

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
	"webapp/globalvar"
)

var (
	LenForCrItems int
	CrNames       []string
)

func GetLen(x []string) map[string]int {
	// execution time
	log.Println("Func GetLen started")
	start := time.Now()
	// declare map
	// <cluster role name>: <len for all their items>
	M1 := make(map[string]int)

	// iterate over slice with cluster role names
	for _, crname := range x {

		//log.Println("Func Clientk8s started ")
		listCr, err := globalvar.Clientset.RbacV1().ClusterRoles().Get(context.TODO(), crname, v1.GetOptions{})
		if err != nil {
			log.Printf("Failed %s", listCr)
			log.Println(err)
		}
		// iterate over cluster role rules
		for _, el := range listCr.Rules {
			// temporary slice for rules
			tempslice := [][]string{el.APIGroups, el.ResourceNames, el.Resources, el.Verbs, el.NonResourceURLs}
			// iterate over tempslice
			for _, y := range tempslice {
				// iterate over items for example APIGroups or ResourceNames then get sum of all len items
				for _, z := range y {
					LenForCrItems += len(z)
				}
			}

		}
		// put name cluster role and len for items to map
		M1[listCr.Name] = LenForCrItems
		// set len to nil, need for the next cluster role items
		LenForCrItems = 0

	}

	/*
		// ##############################################################################################################
		// THIS PART NEED ONLY FOR SORT MAP BY VALUE
		// JUST FOR PRETTY OUTPUT
		pairs := make([][2]interface{}, 0, len(M1))
		for k, v := range M1 {
			pairs = append(pairs, [2]interface{}{k, v})
		}

		// Sort slice based on values
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i][1].(int) < pairs[j][1].(int)
		})

		// Extract sorted keys
		keys := make([]string, len(pairs))
		for i, p := range pairs {
			keys[i] = p[0].(string)
		}

		// Print sorted map
		for _, k := range keys {
			log.Printf("%s: %d\n", k, M1[k])
		}

	*/
	// ##############################################################################################################
	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func GetLen  %s", duration)
	// return map to main
	return M1

}

/*
If we have such cluster role and rules like below, so the len for all items will be 16
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - delete


log.Printf("The len for string namespaces is %d", len("namespaces"))
log.Printf("The len for string delete is %d", len("delete"))
2024/10/22 13:26:13 The len for string namespaces is 10
2024/10/22 13:26:13 The len for string delete is 6
Note: "" won't count
So, after count all items we will know the size(len) for each cluster role and compare it with size(len) provided cluster role


panic: assignment to entry in nil map
You have to initialize the map using the make function (or a map literal) before you can add any elements:

m := make(map[string]float64)
m["pi"] = 3.1416
*/
