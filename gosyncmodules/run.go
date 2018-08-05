package gosyncmodules

import (
	"reflect"
	//"fmt"
	"gopkg.in/ini.v1"
	"gopkg.in/ldap.v2"
	"strings"
	"regexp"
)

func InitialrunAD(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string,
	ADPage int, ADConnTimeout int, UseTLS bool, InsecureSkipVerify bool, CRTValidFor, ADCrtPath string, LDAPBaseDN string,
		shutdownChannel chan string, ADElementsChan chan *[]LDAPElement) {
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
	ADElements := GetFromDS(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	normalizeResult(ADElements, ADBaseDN, LDAPBaseDN)
	//fmt.Println(reflect.TypeOf(ADElements))
	logger.Debugln(ADElements)
	logger.Infoln("Writing results to ", reflect.TypeOf(ADElementsChan))
	logger.Infoln("Length of ", reflect.TypeOf(ADElementsChan), "is", len(ADElements))
	ADElementsChan <- &ADElements
	logger.Infoln("Passing", reflect.TypeOf(ADElementsChan), "to Main thread")

}
func normalizeResult(elements []LDAPElement, ADBasedn, ldapbasedn string)  {
	logger.Debugln("Un-normalized elements from AD:", elements)

	r := regexp.MustCompile(`[A-Z]+=`)


	for i := range elements {
		//elements[i].DN = strings.Replace(strings.ToLower(ADElements[i].DN), ADBaseDN, LDAPBaseDN, -1)
		elements[i].DN = r.ReplaceAllStringFunc(elements[i].DN, func(m string) string {
			return strings.ToLower(m)
		})
		elements[i].DN = strings.Replace(elements[i].DN, ADBasedn, ldapbasedn, -1)
		for i1 := range elements[i].attributes {
			for k, v := range elements[i].attributes[i1] {
				normalizedattribute := ConvertAttributesToLower(&v)
				// replace each element in normalizedattribute with ldapbasedn
				var attributes []string
				for _, attr := range *normalizedattribute {
					a := strings.Replace(attr, ADBasedn, ldapbasedn, -1)
					attributes = append(attributes, a)
				}
				elements[i].attributes[i1][k] = attributes
			}
		}
	}
	logger.Debugln("Normalized AD elements that match the LDAP BaseDN:", elements)

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
	LDAPElements := GetFromDS(connectLDAP, LDAPBaseDN, LDAPFilter, LDAPAttribute, uint32(LDAPPage))
	//Comment below to log ldapelements
	logger.Debugln(LDAPElements)
	logger.Infoln("Length of ", reflect.TypeOf(LDAPElementsChan), "is", len(LDAPElements))

	LDAPElementsChan <- &LDAPElements
	LdapConnectionChan <- connectLDAP

}
