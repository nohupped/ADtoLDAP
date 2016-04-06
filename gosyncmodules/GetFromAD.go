package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
)

func GetFromAD(connect *ldap.Conn, ADBaseDN, ADFilter string, ADAttribute []string, ADPage uint32) *[]ADElement {
	//sizelimit in searchrequest is the limit, which throws an error when the number of results exceeds the limit.
	searchRequest := ldap.NewSearchRequest(ADBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, ADFilter, ADAttribute, nil)
	sr, err := connect.SearchWithPaging(searchRequest, ADPage)
	CheckForError(err)
	fmt.Println(len(sr.Entries))
	ADElements := []ADElement{}
	for _, entry := range sr.Entries{
		NewADEntity := new(ADElement)
		NewADEntity.DN = entry.DN
		for _, attrib := range entry.Attributes {
			NewADEntity.attributes = append(NewADEntity.attributes, keyvalue{attrib.Name: attrib.Values})
		}
		ADElements = append(ADElements, *NewADEntity)
	}
	return &ADElements
}
