package SyncModules

import (
	"github.com/mavricknz/ldap"
	modules "faimodules"
	"time"
)

func ConnectToAD(ADHost string, ldap_port uint16, ADUsername string, ADPassword string)  (*ldap.LDAPConnection){
	modules.Info.Println("Connecting to ldap server", ADHost)
	connection := ldap.NewLDAPConnection(ADHost, ldap_port)
	connection.NetworkConnectTimeout = 10 * time.Second
	modules.Info.Println("Set ldap connect timeout to", connection.NetworkConnectTimeout.Seconds(), "seconds")
	//Connect
	err := connection.Connect()
	modules.CheckForError(err)
	//Bind
	err = connection.Bind(ADUsername, ADPassword)
	//Use the below closure to write to a datastructure using goroutine
	modules.CheckForError(err)
	modules.Info.Println("Binded")
	return connection
}

func GetFromAD(connection *ldap.LDAPConnection, basedn string, filter string, attributes []string, page uint32) []ADElement {
	search_request := ldap.NewSearchRequest(basedn, ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0,false, filter, attributes, nil)
	sr, err := connection.SearchWithPaging(search_request, page)
	modules.CheckForError(err)
	modules.Info.Println("Total elements in AD retrieved: ", len(sr.Entries))
	ADElements := []ADElement{}
	for _, entry := range sr.Entries {
		NewADEntity := new(ADElement)
		NewADEntity.DN = entry.DN
		for _, attrib := range entry.Attributes {
			NewADEntity.attributes = append(NewADEntity.attributes, keyvalue{attrib.Name: attrib.Values})

		}
		ADElements = append(ADElements, *NewADEntity)
	}

	modules.Info.Println(ADElements)
	modules.Info.Println("AD Result gathered")
	return ADElements
}
