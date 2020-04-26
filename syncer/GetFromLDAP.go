package syncer

import (
	"gopkg.in/ldap.v2"
)

// GetFromLDAP retrives values from LDAP / Slave. 
func GetFromLDAP(connect *ldap.Conn, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string, LDAPPage uint32) *[]LDAPElement {
	searchRequest := ldap.NewSearchRequest(LDAPBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, LDAPFilter, LDAPAttribute, nil)
	sr, err := connect.SearchWithPaging(searchRequest, LDAPPage)
	CheckForError(err)
	//fmt.Println(len(sr.Entries))
	ADElements := []LDAPElement{}
	for _, entry := range sr.Entries {
		NewADEntity := new(LDAPElement)
		NewADEntity.DN = entry.DN
		for _, attrib := range entry.Attributes {
			NewADEntity.attributes = append(NewADEntity.attributes, keyvalue{attrib.Name: attrib.Values})
		}
		ADElements = append(ADElements, *NewADEntity)
	}
	return &ADElements
}
