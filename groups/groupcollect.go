// Package groups used only for openshift cluster
package groups

import (
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/utils/strings/slices"
	"log"
	"time"
	"webapp/globalvar"
)

type GroupStruct struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Labels            struct {
				OpenshiftIoLdapHost string `json:"openshift.io/ldap.host"`
			} `json:"labels"`
			Annotations struct {
				OpenshiftIoLdapSyncTime time.Time `json:"openshift.io/ldap.sync-time"`
				OpenshiftIoLdapUID      string    `json:"openshift.io/ldap.uid"`
				OpenshiftIoLdapURL      string    `json:"openshift.io/ldap.url"`
			} `json:"annotations"`
			ManagedFields []struct {
				Manager    string    `json:"manager"`
				Operation  string    `json:"operation"`
				APIVersion string    `json:"apiVersion"`
				Time       time.Time `json:"time"`
				FieldsType string    `json:"fieldsType"`
				FieldsV1   struct {
					FMetadata struct {
						FAnnotations struct {
							NAMING_FAILED struct {
							} `json:"."`
							FOpenshiftIoLdapSyncTime struct {
							} `json:"f:openshift.io/ldap.sync-time"`
							FOpenshiftIoLdapUID struct {
							} `json:"f:openshift.io/ldap.uid"`
							FOpenshiftIoLdapURL struct {
							} `json:"f:openshift.io/ldap.url"`
						} `json:"f:annotations"`
						FLabels struct {
							NAMING_FAILED struct {
							} `json:"."`
							FOpenshiftIoLdapHost struct {
							} `json:"f:openshift.io/ldap.host"`
						} `json:"f:labels"`
					} `json:"f:metadata"`
					FUsers struct {
					} `json:"f:users"`
				} `json:"fieldsV1"`
			} `json:"managedFields"`
		} `json:"metadata"`
		Users []string `json:"users"`
	} `json:"items"`
}

func GroupCollect(LoggedUser string) map[string][]string {

	// get groups with RestClient
	listgroups, err := globalvar.Clientset.AppsV1().RESTClient().Get().AbsPath("/apis/user.openshift.io/v1/groups").DoRaw(context.TODO())
	if err != nil {
		log.Printf("Failed %s", listgroups)
		log.Println(err)
	}

	// init struct
	dataObjet := GroupStruct{}

	jsonErr := json.Unmarshal(listgroups, &dataObjet)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println(dataObjet.Items)

	var SliceGroupForUser []string

	for _, x := range dataObjet.Items {
		log.Println(x.Metadata.Name) // list group name
		log.Println(x.Users)         // list of slice users
		if slices.Contains(x.Users, LoggedUser) {
			SliceGroupForUser = append(SliceGroupForUser, x.Metadata.Name)
		}
	}

	// logged user and them groups
	M1 := map[string][]string{
		LoggedUser: SliceGroupForUser,
	}
	log.Printf("Collectted data for user: %s", M1)
	return M1
}
