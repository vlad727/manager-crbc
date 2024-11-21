// Package getcrname collect all cluster roles names and put it to slice sl, then return slice to main package
package getcrname

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
	"webapp/globalvar"
)

var (
	sl []string
)

// GetCrNameList collect cluster role names and return it to crcheck.go
func GetCrNameList() []string {
	start := time.Now()
	log.Println("Func GetCrNameList started ")
	listCr, err := globalvar.Clientset.RbacV1().ClusterRoles().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("Failed %s", listCr)
		log.Println(err)
	}
	// iterate over cluster roles and append it to slice
	for _, cr := range listCr.Items {
		//log.Println(cr.Name)
		sl = append(sl, cr.Name)
	}
	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func GetCrNameList  %s", duration)
	return sl
}
