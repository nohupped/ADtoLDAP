package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"time"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
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

func ConnectToADTLS(ADHost, AD_Port string, ADUsername, ADPassword string, ADConnTimeout int, CRTInsecureSkipVerify bool,
CRTValidFor, CRTPath string) (*ldap.Conn) {
	ldap.DefaultTimeout = time.Duration(ADConnTimeout) * time.Second
	Info.Println("Set AD connection timeout to", ldap.DefaultTimeout)
	caCert, err := ioutil.ReadFile(CRTPath)
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)


	tlsconfig := new(tls.Config)
	tlsconfig.InsecureSkipVerify = CRTInsecureSkipVerify
	tlsconfig.ServerName = CRTValidFor
	tlsconfig.RootCAs = pool
	Info.Printf("Dialling TLS with config %+v\n", *tlsconfig)
	Info.Println("Nested structs like tls.Config.RootCAs are not logged.")
	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", ADHost, AD_Port), tlsconfig)
	CheckForError(err)
	err = l.Bind(ADUsername, ADPassword)
	CheckForError(err)
	return l
}
