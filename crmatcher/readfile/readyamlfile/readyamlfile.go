// Package readfiles contain struct for yaml cluster role, read yaml file with cluster role
package readyamlfile

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type CrYaml struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		CreationTimestamp time.Time `yaml:"creationTimestamp"`
		Labels            struct {
			ArgocdArgoprojIoInstance string `yaml:"argocd.argoproj.io/instance"`
		} `yaml:"labels"`
		Name            string `yaml:"name"`
		ResourceVersion string `yaml:"resourceVersion"`
		UID             string `yaml:"uid"`
	} `yaml:"metadata"`
	Rules []struct {
		APIGroups       []string `yaml:"apiGroups"`
		ResourceNames   []string `yaml:"resourceNames"`
		Resources       []string `yaml:"resources"`
		Verbs           []string `yaml:"verbs"`
		NonResourceURLs []string `yaml:"nonResourceURLs"`
	} `yaml:"rules"`
}

var (
	LenForCr int
	FileName string
	Cr       CrYaml
)

func ReadFileYaml(dstDir string) {

	// logging readFile
	log.Println("Func ReadFileYaml started")

	// read yaml file
	yamlFile, err := os.ReadFile(dstDir)
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}
	err = yaml.Unmarshal(yamlFile, &Cr)
	if err != nil {
		panic(err)
	}

	// get len for cluster role from provided yaml
	for _, el := range Cr.Rules {
		//log.Println(el)
		tempslice := [][]string{el.APIGroups, el.ResourceNames, el.Resources, el.Verbs, el.NonResourceURLs}
		for _, y := range tempslice {
			for _, x := range y {
				LenForCr += len(x)
			}
		}
	}
	log.Printf("The len for %s Cluster Role is %d", Cr.Metadata.Name, LenForCr)

}
