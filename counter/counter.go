// Package counter will set label for cluster role binding with end "crbc" get int and export it to package crbcmain
// func Counter started in func main
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

func Counter() {

	// // Code to measure
	start := time.Now()

	// logging
	log.Println("Func Counter started")

	listCRB, err := globalvar.Clientset.RbacV1().ClusterRoleBindings().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		log.Println("Can't list cluster role bindings")
	}

	for _, x := range listCRB.Items {
		//log.Println(x.Name)
		if strings.Contains(x.Name, "crbc") && x.Name != "manager-crbc-clusterrolebinding-admin" {
			CreatedByCrbc = append(CreatedByCrbc, x.Name)
		}
	}
	log.Println(CreatedByCrbc)
	NumberOfEntities = len(CreatedByCrbc)
	log.Printf("Cluster contain %d cluster role bindings created by manager-crbc", len(CreatedByCrbc))

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func Counter  %s", duration)

}
