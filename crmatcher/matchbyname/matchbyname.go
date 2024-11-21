// Package matchbyname compare cluster role name from yaml and allowed cluster role from slice
package matchbyname

import "k8s.io/utils/strings/slices"

func MatchByName(x []string, y string) bool {

	value := slices.Contains(x, y)
	return value

}
