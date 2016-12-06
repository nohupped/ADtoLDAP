package gosyncmodules

import (
	"reflect"
	//"fmt"
	"gopkg.in/ini.v1"
	"gopkg.in/ldap.v2"
)

func InitialrunAD(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string,
	ADPage int, ADConnTimeout int, UseTLS bool, InsecureSkipVerify bool, CRTValidFor, ADCrtPath string, shutdownChannel chan string, ADElementsChan chan *[]LDAPElement)  {
	Info.Println("Connecting to AD", ADHost)
	var connectAD *ldap.Conn
	if UseTLS == false {
		connectAD = ConnectToDirectoryServer(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout)
	} else  {
		connectAD = ConnectToDirectoryServerTLS(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout, InsecureSkipVerify, CRTValidFor, ADCrtPath)
	}
	defer func() {shutdownChannel <- "Done from func InitialrunAD"}()
	defer Info.Println("closed")
	defer connectAD.Close()
	defer Info.Println("Closing connection")
	ADElements := GetFromAD(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	//fmt.Println(reflect.TypeOf(ADElements))
//	Info.Println(ADElements)
	Info.Println("Writing results to ", reflect.TypeOf(ADElementsChan))
	Info.Println("Length of ", reflect.TypeOf(ADElementsChan), "is", len(*ADElements))
	ADElementsChan <- ADElements
	Info.Println("Passing", reflect.TypeOf(ADElementsChan), "to Main thread")


}

func InitialrunLDAP(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string,
	LDAPPage int, LDAPConnTimeout int, ADElements *[]LDAPElement, ReplaceAttributes, MapAttributes *ini.Section)  {
	Info.Println("Received", len(*ADElements), "elements to populate ldap")
	connectLDAP := ConnectToDirectoryServer(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	InitialPopulateToLdap(ADElements, connectLDAP, ReplaceAttributes, MapAttributes, false)
	defer Info.Println("closed")
	defer connectLDAP.Close()
	defer Info.Println("Closing connection")
	Info.Println("End of LDAP")

}

func SyncrunLDAP(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string,
			LDAPPage int, LDAPConnTimeout int, shutdownChannel chan string,
			LDAPElementsChan chan *[]LDAPElement, LdapConnectionChan chan *ldap.Conn,
			ReplaceAttributes, MapAttributes *ini.Section)  {
	Info.Println("Connecting to LDAP", LDAPHost)
	connectLDAP := ConnectToDirectoryServer(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	defer func() {shutdownChannel <- "Done from func syncrunLDAP"}()
	LDAPElements := GetFromLDAP(connectLDAP, LDAPBaseDN, LDAPFilter, LDAPAttribute, uint32(LDAPPage))
	//Comment below to log ldapelements
	//Info.Println(LDAPElements)
	Info.Println("Length of ", reflect.TypeOf(LDAPElementsChan), "is", len(*LDAPElements))

	LDAPElementsChan <- LDAPElements
	LdapConnectionChan <- connectLDAP

}
