package trimmer

import (
	"log"
	"strings"
)

func Trimmer(x map[string][]string) string {
	// temp var
	var tmp string
	// itearte over slice
	for k, v := range x {
		log.Printf("Key: %s Value: %s", k, v)
		// slice to string
		strtoken := strings.Join(v, " ")
		// delete unused part of token
		tmp = strings.ReplaceAll(strtoken, "Bearer ", "")
	}
	return tmp
}
