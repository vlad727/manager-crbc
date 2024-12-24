// Package counter count number of cluster role binding with end "crbc" get int and return it to package crbcmain
package counter

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
	"time"
	"webapp/globalvar"
)

var (
	CreatedByCrbc    []string
	NumberOfEntities int
)

func Counter() int {

	// // Code to measure
	start := time.Now()

	// logging
	log.Println("Func Counter started")

	// get all cluster role bindings
	listCRB, err := globalvar.Clientset.RbacV1().ClusterRoleBindings().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		log.Println("Can't list cluster role bindings")
	}
	// map of cluster role bindings
	countCrb := make(map[string]string)
	for _, x := range listCRB.Items {
		if strings.Contains(x.Name, "crbc") && x.Name != "manager-crbc-clusterrolebinding-admin" {
			countCrb[x.Name] = x.Kind

		}
	}
	// logging
	log.Println(CreatedByCrbc)
	NumberOfEntities = len(countCrb)
	log.Printf("Cluster contain %d cluster role bindings created by manager-crbc", len(CreatedByCrbc))

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func Counter  %s", duration)
	return NumberOfEntities

}
