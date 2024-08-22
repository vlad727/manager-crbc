package validate_bac

import (
	"encoding/json"
	"k8s.io/api/admission/v1beta1"
	"log"
	"net/http"
)

var (
	RequesterForCrb string
)

// Validate handler accepts or rejects based on request contents
func Validate(w http.ResponseWriter, r *http.Request) {

	// logging
	log.Println("Validating request...")
	log.Println("Func validate.bac ...")
	log.Println("Waiting cluster role binding...")

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

	// admission request get requester name
	RequesterForCrb = arReview.Request.UserInfo.Username
	log.Printf("Requester for crb is %s", RequesterForCrb)
	// get requested cluster role binding name
	RequestedCrbName := arReview.Request.Name //2024/07/30 15:25:30 subresName4 cilium-operator-kube-system-admin-cluster-role-crbc
	log.Printf("Requested cluster role binding name %s", RequestedCrbName)

	// response
	arReview.Response = &v1beta1.AdmissionResponse{
		UID:     arReview.Request.UID,
		Allowed: true,
	}
	//log.Println(arReview.Response)
	// 2024/07/02 11:46:55 &AdmissionResponse{UID:f12d6085-2aa5-476b-92b4-a8d3ba9d219e,Allowed:true,Result:nil,Patch:nil,PatchType:nil,AuditAnnotations:map[string]string{},Warnings:[],}
	//log.Println("The end of func validate.bac")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&arReview)
}
