package gosyncmodules

import (
	"github.com/nohupped/ldap"
)

type Action map[string]*ldap.AddRequest

func FindAdds(ADElementsConverted, LDAPElementsConverted *[]*ldap.AddRequest, AddChan chan Action, shutdownChannel chan string){
	Info.Println("Starting FindAdds")
	defer func() {shutdownChannel <- "Done from func FindAdds"}()
	defer close(AddChan)
	defer Info.Println("About to close blocking channel from FindAdds")
	for _, i := range *ADElementsConverted {
		if IfDNExists(i, *LDAPElementsConverted) {
			Info.Println(i, "exists, checking for change in attributes.")
			continue
		} else {
			//err := LDAPConnection.Add(i)
			AddChan <- Action{"Add":i}

		}
	}
}