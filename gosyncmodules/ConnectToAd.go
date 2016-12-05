package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"time"
	"crypto/tls"
)

func ConnectToAD(ADHost, AD_Port string, ADUsername, ADPassword string, ADConnTimeout int) (*ldap.Conn){
	ldap.DefaultTimeout = time.Duration(ADConnTimeout) * time.Second
	Info.Println("Set AD connection timeout to", ldap.DefaultTimeout)
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ADHost, AD_Port))
	CheckForError(err)
	Info.Println("Binding")
	err = l.Bind(ADUsername, ADPassword)
	CheckForError(err)
	return l
}

func ConnectToADTLS(ADHost, AD_Port string, ADUsername, ADPassword string, ADConnTimeout int, CRTInsecureSkipVerify bool, CRTValidFor string) (*ldap.Conn) {
	ldap.DefaultTimeout = time.Duration(ADConnTimeout) * time.Second
	Info.Println("Set AD connection timeout to", ldap.DefaultTimeout)
	Info.Println("Dialling TLS")
	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", ADHost, AD_Port), &tls.Config{InsecureSkipVerify: CRTInsecureSkipVerify, ServerName: CRTValidFor})
	CheckForError(err)
	err = l.Bind(ADUsername, ADPassword)
	CheckForError(err)
	return l
}
