package gosyncmodules

import (
	"github.com/nohupped/ldap"
)

type Action map[string]*ldap.AddRequest

func FindAdds(ADElementsConverted, LDAPElementsConverted *[]*ldap.AddRequest, LdapConnection *ldap.Conn, AddChan chan Action, shutdownAddChan chan string){
	Info.Println("Starting FindAdds")
	defer func() {shutdownAddChan <- "Done from func FindAdds"}()
	defer close(AddChan)
	defer Info.Println("About to close blocking channel from FindAdds")
	for _, i := range *ADElementsConverted {
		Exists, LDAPEntry := IfDNExists(i, *LDAPElementsConverted)
		if Exists {
			Info.Println(i, "exists, checking for change in attributes.")
			CheckAttributes(LdapConnection, LDAPEntry, i)
			continue
		} else {
			//err := LDAPConnection.Add(i)
			AddChan <- Action{"Add":i}  //Write composite literal to channel

		}
	}
}


func FindDels(LDAPElementsConverted, ADElementsConverted *[]*ldap.AddRequest, DelChan chan Action, shutdownDelChan chan string){
	Info.Println("Starting FindDels")
	defer func() {shutdownDelChan <- "Done from func FindDels"}()
	defer close(DelChan)
	defer Info.Println("About to close blocking channel from FindAdds")
	for _, i := range *LDAPElementsConverted {
		Exists, _ := IfDNExists(i, *ADElementsConverted)
		if Exists {
			continue
		} else {
			Info.Println(i.DN, "Doesn't exist in AD, will be set to delete.")
			//err := LDAPConnection.Add(i)
			DelChan <- Action{"Del":i}  //Write composite literal to channel

		}
	}
}

