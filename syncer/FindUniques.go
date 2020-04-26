package syncer

import (
	//	"github.com/nohupped/ldap"
	"gopkg.in/ldap.v2"
)

// Action map is a lame attempt to define if the action needs to be an Add or a Delete
type Action map[string]*ldap.AddRequest

// FindAdds function is used to find the entries that needs to be modified/added to the slave
func FindAdds(ADElementsConverted, LDAPElementsConverted *[]*ldap.AddRequest, LdapConnection *ldap.Conn, AddChan chan Action, shutdownAddChan chan string) {
	logger.Debugln("Starting FindAdds")
	defer func() { shutdownAddChan <- "Done from func FindAdds" }()
	defer close(AddChan)
	defer logger.Debugln("About to close blocking channel from FindAdds")
	for _, i := range *ADElementsConverted {
		Exists, LDAPEntry := IfDNExists(i, *LDAPElementsConverted)
		if Exists {
			logger.Debugln(i, "exists, checking for change in attributes.")
			CheckAttributes(LdapConnection, LDAPEntry, i)
			continue
		} else {
			AddChan <- Action{"Add": i} //Write composite literal to channel

		}
	}
}

// FindDels function is used to find the entries that needs to be deleted from the slave
func FindDels(LDAPElementsConverted, ADElementsConverted *[]*ldap.AddRequest, DelChan chan Action, shutdownDelChan chan string) {
	logger.Debugln("Starting FindDels")
	defer func() { shutdownDelChan <- "Done from func FindDels" }()
	defer close(DelChan)
	defer logger.Debugln("About to close blocking channel from FindAdds")
	for _, i := range *LDAPElementsConverted {
		Exists, _ := IfDNExists(i, *ADElementsConverted)
		if Exists {
			continue
		} else {
			logger.Debugln(i.DN, "Doesn't exist in AD, will be set to delete.")
			DelChan <- Action{"Del": i} //Write composite literal to channel

		}
	}
}
