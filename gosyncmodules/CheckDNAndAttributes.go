package gosyncmodules

import (
	"github.com/nohupped/ldap" //using a forked version that includes custom methods to retrieve and edit *AddRequest struct.
	//"fmt"
	"fmt"
)

func IfDNExists(checkfor *ldap.AddRequest , checkin []*ldap.AddRequest ) (bool, *ldap.AddRequest) {
	for _, i := range checkin {
		if checkfor.DN == i.DN {
			return true, i
		}

	}
	return false, nil
}

func CheckAttributes(LdapEntry, ADEntry *ldap.AddRequest)  {

	for _, adEntries := range ADEntry.Attributes {
		fmt.Println(adEntries.AttrType, "::", adEntries.AttrVals)
	}


}



