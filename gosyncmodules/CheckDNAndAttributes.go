package gosyncmodules

import (
	"gopkg.in/ldap.v2"
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
		if adEntries.Type == "memberOf" {
			adEntries.Vals = *ConvertAttributesToLower(&adEntries.Vals)
		}
		sort.Strings(adEntries.Vals)
		ADMapped := MapADandLDAP{adEntries.Type: adEntries.Vals}
		ADMapAggregated = append(ADMapAggregated, ADMapped)
	}
	for _, ldapEntries  := range LdapEntry.Attributes {
		sort.Strings(ldapEntries.Vals)
		LDAPMapped := MapADandLDAP{ldapEntries.Type: ldapEntries.Vals}
		LDAPMapAggregated = append(LDAPMapAggregated, LDAPMapped)
	}

	logger.Debugln("Got from AD", ADMapAggregated)
	logger.Debugln("Got from LD", LDAPMapAggregated)

	if reflect.DeepEqual(ADMapAggregated, LDAPMapAggregated) == true {
		logger.Debugln("Both entries matches, passing...")
	} else {
		logger.Debugln("CHANGE DETECTED")
		logger.Debugln("AD -> ", ADMapAggregated)
		logger.Debugln("LD -> ", LDAPMapAggregated)
		delete := ldap.NewDelRequest(LdapEntry.DN, []ldap.Control{})
		err := LdapConnection.Del(delete)
		if err != nil {
			logger.Errorln(err)
		} else {logger.Debugln(*delete, "Deleted")}
		err = LdapConnection.Add(ADEntry)
		if err != nil {
			logger.Debugln(err)
		} else {logger.Debugln(*ADEntry, "Added to ldap")}

	}


}



