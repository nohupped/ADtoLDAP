package gosyncmodules

import (
	"github.com/nohupped/ldap" //using a forked version that includes custom methods to retrieve and edit *AddRequest struct.
	//"fmt"
)

func GetFromAD(connect *ldap.Conn, ADBaseDN, ADFilter string, ADAttribute []string, ADPage uint32) *[]LDAPElement {
	//sizelimit in searchrequest is the limit, which throws an error when the number of results exceeds the limit.
	searchRequest := ldap.NewSearchRequest(ADBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, ADFilter, ADAttribute, nil)
	sr, err := connect.SearchWithPaging(searchRequest, ADPage)
	CheckForError(err)
	//fmt.Println(len(sr.Entries))
	ADElements := []LDAPElement{}
	for _, entry := range sr.Entries{
		NewADEntity := new(LDAPElement)
		NewADEntity.DN = entry.DN
		for _, attrib := range entry.Attributes {
			NewADEntity.attributes = append(NewADEntity.attributes, keyvalue{attrib.Name: attrib.Values})
		}
		ADElements = append(ADElements, *NewADEntity)
	}
	return &ADElements
}
