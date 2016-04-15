package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
)

func IfDuplicateExists(checkfor *ldap.AddRequest , checkin []*ldap.AddRequest ) bool {

	for _, i := range checkin {
		fmt.Println("Checking if ", *checkfor, "is equal to ", *i)
		if i == checkfor {
			return true
		}
	}
	return false
}
