package gosyncmodules

import (
	"reflect"
	//"fmt"
	"gopkg.in/ini.v1"
	"gopkg.in/ldap.v2"
)

func InitialrunAD(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string,
	ADPage int, ADConnTimeout int, UseTLS bool, InsecureSkipVerify bool, CRTValidFor, ADCrtPath string, shutdownChannel chan string, ADElementsChan chan *[]LDAPElement) {
	logger.Infoln("Connecting to AD", ADHost)
	var connectAD *ldap.Conn
	if UseTLS == false {
		connectAD = ConnectToDirectoryServer(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout)
	} else {
		connectAD = ConnectToDirectoryServerTLS(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout, InsecureSkipVerify, CRTValidFor, ADCrtPath)
	}
	defer func() { shutdownChannel <- "Done from func InitialrunAD" }()
	defer logger.Infoln("closed")
	defer connectAD.Close()
	defer logger.Infoln("Closing connection")
	ADElements := GetFromAD(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	//fmt.Println(reflect.TypeOf(ADElements))
	//	Info.Println(ADElements)
	logger.Infoln("Writing results to ", reflect.TypeOf(ADElementsChan))
	logger.Infoln("Length of ", reflect.TypeOf(ADElementsChan), "is", len(*ADElements))
	ADElementsChan <- ADElements
	logger.Infoln("Passing", reflect.TypeOf(ADElementsChan), "to Main thread")

}

func InitialrunLDAP(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string,
	LDAPPage int, LDAPConnTimeout int, LDAPUseTLS bool, LDAPCrtValidFor string, LDAPCrtPath string, LDAPCRTInsecureSkipVerify bool,
	ADElements *[]LDAPElement, ReplaceAttributes, MapAttributes *ini.Section) {
	logger.Infoln("Received", len(*ADElements), "elements to populate ldap")
	var connectLDAP *ldap.Conn
	if LDAPUseTLS == false {
		connectLDAP = ConnectToDirectoryServer(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	} else {
		connectLDAP = ConnectToDirectoryServerTLS(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout,
			LDAPCRTInsecureSkipVerify, LDAPCrtValidFor, LDAPCrtPath)
	}

	InitialPopulateToLdap(ADElements, connectLDAP, ReplaceAttributes, MapAttributes, false)
	defer logger.Infoln("closed")
	defer connectLDAP.Close()
	defer logger.Infoln("Closing connection")
	logger.Infoln("End of LDAP")

}

func SyncrunLDAP(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string,
	LDAPPage int, LDAPConnTimeout int, LDAPUseTLS bool, LDAPCRTInsecureSkipVerify bool, LDAPCrtValidFor string,
	LDAPCrtPath string, shutdownChannel chan string, LDAPElementsChan chan *[]LDAPElement,
	LdapConnectionChan chan *ldap.Conn, ReplaceAttributes, MapAttributes *ini.Section) {
	logger.Infoln("Connecting to LDAP", LDAPHost)
	//connectLDAP := ConnectToDirectoryServer(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	var connectLDAP *ldap.Conn
	if LDAPUseTLS == false {
		connectLDAP = ConnectToDirectoryServer(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	} else {
		connectLDAP = ConnectToDirectoryServerTLS(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout,
			LDAPCRTInsecureSkipVerify, LDAPCrtValidFor, LDAPCrtPath)
	}
	defer func() { shutdownChannel <- "Done from func syncrunLDAP" }()
	LDAPElements := GetFromLDAP(connectLDAP, LDAPBaseDN, LDAPFilter, LDAPAttribute, uint32(LDAPPage))
	//Comment below to log ldapelements
	//Info.Println(LDAPElements)
	logger.Infoln("Length of ", reflect.TypeOf(LDAPElementsChan), "is", len(*LDAPElements))

	LDAPElementsChan <- LDAPElements
	LdapConnectionChan <- connectLDAP

}
