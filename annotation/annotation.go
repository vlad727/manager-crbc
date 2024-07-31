package annotation

import (
	"encoding/json"
	"golang.org/x/net/context"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
	"time"
	"webapp/globalvar"
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

// Validate handler accepts or rejects based on request contents
func Validate(w http.ResponseWriter, r *http.Request) {

	log.Println("Validating request...")
	log.Println("Func validate ...")

	// var arReview with struct v1beta1.AdmissionReview{}
	arReview := v1beta1.AdmissionReview{}

	// decode arReview to json and check request
	if err := json.NewDecoder(r.Body).Decode(&arReview); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if arReview.Request == nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// admission request get crb name
	requestedCrbName := arReview.Request.Name //2024/07/30 15:25:30 subresName4 cilium-operator-kube-system-admin-cluster-role-crbc
	log.Printf("Requested clsuter role binding name %s", requestedCrbName)

	// like example for annotation created-by: crbc-manager
	// crb-requester: <ldap-user>
	// wait while namespace will created
	time.Sleep(1 * time.Second)
	setAnnotation := mainstruct{
		Metadata: Annotations{
			Requester{"ldap-user"},
		},
	}

	// marshal var setAnnotation to json
	bytes, _ := json.Marshal(setAnnotation)

	// set annotation to cluster role binding
	//Note: that type used MergePatchType (allow add new piece of json)
	_, err := globalvar.Clientset.RbacV1().ClusterRoleBindings().Patch(context.TODO(), requestedCrbName, types.MergePatchType, bytes, metav1.PatchOptions{})
	if err != nil {
		log.Printf("Failed to set annotation for %s", requestedCrbName)
		log.Println(err)
	} else {
		log.Println("Namespace has been annotated with", string(bytes))
	}

	arReview.Response = &v1beta1.AdmissionResponse{
		UID:     arReview.Request.UID,
		Allowed: true,
	}
	//log.Println(arReview.Response)
	// 2024/07/02 11:46:55 &AdmissionResponse{UID:f12d6085-2aa5-476b-92b4-a8d3ba9d219e,Allowed:true,Result:nil,Patch:nil,PatchType:nil,AuditAnnotations:map[string]string{},Warnings:[],}
	//log.Println("The end of func validate")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&arReview)
}
