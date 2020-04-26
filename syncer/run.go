package syncer

import (
	"reflect"
	"gopkg.in/ini.v1"
	"gopkg.in/ldap.v2"
)

// SyncRunAD connects to AD for sync and writes the retrieved AD elements to the channel.
func SyncRunAD(ADHost, ADPort, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string,
	ADPage int, ADConnTimeout int, UseTLS bool, InsecureSkipVerify bool, CRTValidFor, ADCrtPath string, shutdownChannel chan string, ADElementsChan chan *[]LDAPElement) {
	logger.Infoln("Connecting to AD", ADHost)
	var connectAD *ldap.Conn
	if UseTLS == false {
		connectAD = ConnectToDirectoryServer(ADHost, ADPort, ADUsername, ADPassword, ADConnTimeout)
	} else {
		connectAD = ConnectToDirectoryServerTLS(ADHost, ADPort, ADUsername, ADPassword, ADConnTimeout, InsecureSkipVerify, CRTValidFor, ADCrtPath)
	}
	defer func() { shutdownChannel <- "Done from func InitialrunAD" }()
	defer logger.Infoln("closed")
	defer connectAD.Close()
	defer logger.Infoln("Closing connection")
	ADElements := GetFromAD(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	logger.Infoln("Writing results to ", reflect.TypeOf(ADElementsChan))
	logger.Infoln("Length of ", reflect.TypeOf(ADElementsChan), "is", len(*ADElements))
	ADElementsChan <- ADElements
	logger.Infoln("Passing", reflect.TypeOf(ADElementsChan), "to Main thread")

}

// SyncRunLDAP connects to LDAP for sync, and writes the processed elements to channel.
func SyncRunLDAP(LDAPHost, LDAPPort, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string,
	LDAPPage int, LDAPConnTimeout int, LDAPUseTLS bool, LDAPCRTInsecureSkipVerify bool, LDAPCrtValidFor string,
	LDAPCrtPath string, shutdownChannel chan string, LDAPElementsChan chan *[]LDAPElement,
	LdapConnectionChan chan *ldap.Conn, ReplaceAttributes, MapAttributes *ini.Section) {
	logger.Infoln("Connecting to LDAP", LDAPHost)
	var connectLDAP *ldap.Conn
	if LDAPUseTLS == false {
		connectLDAP = ConnectToDirectoryServer(LDAPHost, LDAPPort, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	} else {
		connectLDAP = ConnectToDirectoryServerTLS(LDAPHost, LDAPPort, LDAPUsername, LDAPPassword, LDAPConnTimeout,
			LDAPCRTInsecureSkipVerify, LDAPCrtValidFor, LDAPCrtPath)
	}
	defer func() { shutdownChannel <- "Done from func syncrunLDAP" }()
	LDAPElements := GetFromLDAP(connectLDAP, LDAPBaseDN, LDAPFilter, LDAPAttribute, uint32(LDAPPage))

	logger.Infoln("Length of ", reflect.TypeOf(LDAPElementsChan), "is", len(*LDAPElements))

	LDAPElementsChan <- LDAPElements
	LdapConnectionChan <- connectLDAP

}
