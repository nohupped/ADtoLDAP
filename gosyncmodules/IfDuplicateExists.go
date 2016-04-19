package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	//"fmt"
)

func IfDNExists(checkfor *ldap.AddRequest , checkin []*ldap.AddRequest ) bool {
	for _, i := range checkin {
	//	fmt.Println("Checking", checkfor.GetDN(), "equals", i.GetDN())
		if checkfor.GetDN() == i.GetDN() {
			return true
		}

	}
	return false
}
