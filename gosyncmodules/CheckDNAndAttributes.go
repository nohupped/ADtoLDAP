package gosyncmodules

import (
	"github.com/nohupped/ldap" //using a forked version that includes custom methods to retrieve and edit *AddRequest struct.
	"sort"
	"reflect"
)

func IfDNExists(checkfor *ldap.AddRequest , checkin []*ldap.AddRequest ) (bool, *ldap.AddRequest) {
	for _, i := range checkin {
		if checkfor.DN == i.DN {
			return true, i
		}

	}
	return false, nil
}

type MapADandLDAP map[string][]string

func CheckAttributes(LdapConnection *ldap.Conn, LdapEntry, ADEntry *ldap.AddRequest)  {
	var ADMapAggregated []MapADandLDAP
	var LDAPMapAggregated []MapADandLDAP
	for _, adEntries := range ADEntry.Attributes {
		if adEntries.AttrType == "memberOf" {
			adEntries.AttrVals = *ConvertAttributesToLower(&adEntries.AttrVals)
		}
		sort.Strings(adEntries.AttrVals)
		ADMapped := MapADandLDAP{adEntries.AttrType: adEntries.AttrVals}
		ADMapAggregated = append(ADMapAggregated, ADMapped)
	}
	for _, ldapEntries  := range LdapEntry.Attributes {
		sort.Strings(ldapEntries.AttrVals)
		LDAPMapped := MapADandLDAP{ldapEntries.AttrType: ldapEntries.AttrVals}
		LDAPMapAggregated = append(LDAPMapAggregated, LDAPMapped)
	}

	Info.Println("Got from AD", ADMapAggregated, "\n")
	Info.Println("Got from LD", LDAPMapAggregated, "\n\n")

	if reflect.DeepEqual(ADMapAggregated, LDAPMapAggregated) == true {
		Info.Println("Both entries matches, passing...")
	} else {
		Info.Println("CHANGE DETECTED, \nAD ->", ADMapAggregated, "\nLD ->", LDAPMapAggregated, "\n\n")
		delete := ldap.NewDelRequest(LdapEntry.DN, []ldap.Control{})
		err := LdapConnection.Del(delete)
		if err != nil {
			Error.Println(err)
		} else {Info.Println(*delete, "Deleted")}
		err = LdapConnection.Add(ADEntry)
		if err != nil {
			Error.Println(err)
		} else {Info.Println(*ADEntry, "Added to ldap")}

	}


}



