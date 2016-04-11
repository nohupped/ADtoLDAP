package gosyncmodules

import (
	"reflect"
	//"fmt"
	"gopkg.in/ini.v1"
)

func InitialrunAD(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string,
	ADPage int, ADConnTimeout int, shutdownChannel chan string, ADElementsChan chan *[]ADElement)  {
	connectAD := ConnectToAD(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout)
	defer func() {shutdownChannel <- "Done from func InitialrunAD"}()
	defer Info.Println("closed")
	defer connectAD.Close()
	defer Info.Println("Closing connection")
	ADElements := GetFromAD(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	//fmt.Println(reflect.TypeOf(ADElements))
	Info.Println(ADElements)
	Info.Println("Writing results to ", reflect.TypeOf(ADElementsChan))
	Info.Println("Length of ", reflect.TypeOf(ADElementsChan), "is", len(*ADElements))
	ADElementsChan <- ADElements
	Info.Println("Passing", reflect.TypeOf(ADElementsChan), "to Main thread")


}

func InitialrunLDAP(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string,
	LDAPPage int, LDAPConnTimeout int, ADElements *[]ADElement, ReplaceAttributes, MapAttributes *ini.Section)  {
	Info.Println("Received", len(*ADElements), "elements to populate ldap")
	connectLDAP := ConnectToLdap(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	InitialPopulateToLdap(ADElements, connectLDAP, ReplaceAttributes, MapAttributes)
	defer Info.Println("closed")
	defer connectLDAP.Close()
	defer Info.Println("Closing connection")
	Info.Println("End of LDAP")

}
