package parsepost

import (
	"encoding/json"
	"golang.org/x/net/context"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
	"strings"
	"webapp/globalvar"
	"webapp/jwtdecode"
)

var (
	sa       = "" // service account name
	ns       = "" // namespace name (namespace where service account exist)
	cr       = "" // cluster role (requested cluster role which one chose by user)
	Crbname  = "" // cluster role binding name concatenate sa + ns + cl ( var for logging to get result name for crb)
	ErrorMsg = "" // Message for page errormsg (if crb already exist or smt else)
	Checkbox = ""
)

// struct for json
type mainstruct struct {
	Metadata Annotations `json:"metadata"`
}

type Annotations struct {
	Annotations Requester `json:"annotations"`
}
type Requester struct {
	Requester string `json:"requester"`
}

func bindingSubjects(saName, namespace string) []rbacv1.Subject {

	if Checkbox != "" {
		return []rbacv1.Subject{
			{
				Kind:      rbacv1.UserKind,
				Name:      saName,
				Namespace: namespace,
			},
		}
	} else {
		return []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      saName,
				Namespace: namespace,
			},
		}
	}

}

func ParsePostRequest(w http.ResponseWriter, r *http.Request) {

	// ----------------------------------------------------------------------------------------------------------------------------
	// parse post request

	// init empty slice
	sl := []string{}
	r.ParseForm() // Анализирует переданные параметры url, затем анализирует пакет ответа для тела POST (тела запроса)
	// внимание: без вызова метода ParseForm последующие данные не будут получены
	log.Printf("Full post request: %s", r)
	log.Println(r.Form) // печатает информацию на сервере
	log.Println("Path: ", r.URL.Path)
	//log.Println("Schema: ", r.URL.Scheme)
	//log.Println(r.Form["url_long"])
	// iterate over map
	for k, v := range r.Form {
		log.Println("Key: ", k)
		//fmt.Println("Value: ", strings.Join(v, " "))
		log.Println(v)

		// check checkbox
		if k == "CrbLikeUser" {
			log.Println("Need to set \"- kind: User\"")
			Checkbox = "True"

		}
		if k == "choice1" {
			// split string in slice
			for _, el := range v {
				if strings.Contains(el, " ") {
					substrs := strings.Split(el, " ")
					for _, element := range substrs {
						sl = append(sl, element)
					}
				} else {
					sl = append(sl, el)
				}
			}
		}
	}
	// ----------------------------------------------------------------------------------------------------------------------------
	// create cluster role binding logging input

	// sl it's slice with service account namespace and requested cluster role
	for index, el := range sl {
		//log.Println(index, el)
		switch index {
		case 0:
			ns = el
			log.Printf("The namespace is %s", ns)
			// Using the ReplaceAll Function
			resultDelColon := strings.ReplaceAll(ns, ":", "")
			ns = resultDelColon
		case 1:
			sa = el
			log.Printf("The service account is %s", sa)
		case 2:
			cr = el
			log.Printf("The cluster role is %s", cr)

		}
	}

	// init var cluster role binding for service account
	binding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: sa + "-" + ns + "-" + cr + "-" + "crbc",
		},
		Subjects: bindingSubjects(sa, ns),
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     cr,
		},
	}

	// create cluster role binding with clientset
	_, err := globalvar.Clientset.RbacV1().ClusterRoleBindings().Create(context.Background(), binding, v1.CreateOptions{})
	if err != nil {
		log.Println(err)
		ErrorMsg = err.Error()
		// redirect to failed creation page
		// if crb already exist or smt goes wrong
		http.Redirect(w, r, "/error", http.StatusSeeOther)

	} else {
		// concatenate strings to crb name
		Crbname = sa + "-" + ns + "-" + cr + "-" + "crbc"
		// redirect to success creation page
		// show page with crb name
		http.Redirect(w, r, "/crbshow", http.StatusSeeOther)

		log.Printf("Cluster role binding %s has been created...", Crbname)

		// set annotation
		// validate.bac example: crb-requester: <ldap-user>
		setAnnotation := mainstruct{
			Metadata: Annotations{
				Requester{jwtdecode.LoggedUser},
			},
		}

		// marshal var setAnnotation to json
		bytes, _ := json.Marshal(setAnnotation)

		//Note: that type used MergePatchType (allow add new piece of json)
		_, err = globalvar.Clientset.RbacV1().ClusterRoleBindings().Patch(context.TODO(), Crbname, types.MergePatchType, bytes, v1.PatchOptions{})
		if err != nil {
			log.Printf("Failed to set validate.bac for %s", Crbname)
			log.Println(err)
		} else {
			log.Println("Cluster role binding has been annotated", string(bytes))
		}

	}
}
