package gosyncmodules

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"time"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
)

// ConnectToDirectoryServer will try to establish a connection to the directory server and return the connection
// object. This is an un-encrypted connection, and the data transferred will be human readable if checked with tcpdump
// eg: tcpdump -v -XX
func ConnectToDirectoryServer(Host, Port string, Username, Password string, ConnTimeout int) (*ldap.Conn){
	ldap.DefaultTimeout = time.Duration(ConnTimeout) * time.Second
	logger.Debugln("Set AD connection timeout to", ldap.DefaultTimeout)
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", Host, Port))
	CheckForError(err)
	logger.Debugln("Binding")
	err = l.Bind(Username, Password)
	CheckForError(err)
	return l
}

// ConnectToDirectoryServerTLS will try to establish a tls encrypted connection to the directory server, and return
// the connection object. If the CRTInsecureSkipVerify is set to false, this function will read the pem file from CRTPath
// to add the certificates into a cert pool and use it as the Root CAs, and set the ServerName to which
// the certificate was issued for, from CRTValidFor.
func ConnectToDirectoryServerTLS(Host, Port string, Username, Password string, ConnTimeout int, CRTInsecureSkipVerify bool,
CRTValidFor, CRTPath string) (*ldap.Conn) {
	ldap.DefaultTimeout = time.Duration(ConnTimeout) * time.Second
	logger.Debugln("Set AD connection timeout to", ldap.DefaultTimeout)
	tlsconfig := new(tls.Config)

	if CRTInsecureSkipVerify == false {
		tlsconfig = tlsconfigNoSkipVerify(CRTInsecureSkipVerify, CRTValidFor, CRTPath)
	} else {
		tlsconfig = tlsconfigSkipVerify(CRTInsecureSkipVerify)
	}

	logger.Debugf("Dialling TLS with config %+v\n", *tlsconfig)
	logger.Debugf("Nested structs like tls.Config.RootCAs are not logged.")
	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%s", Host, Port), tlsconfig)
	CheckForError(err)
	err = l.Bind(Username, Password)
	CheckForError(err)
	return l
}

func tlsconfigSkipVerify(CRTInsecureSkipVerify bool) *tls.Config {
	tlsconfig := new(tls.Config)
	tlsconfig.InsecureSkipVerify = CRTInsecureSkipVerify
	return tlsconfig
}

func tlsconfigNoSkipVerify(CRTInsecureSkipVerify bool, CRTValidFor, CRTPath string)  *tls.Config{
	tlsconfig := new(tls.Config)
	caCert, err := ioutil.ReadFile(CRTPath)
	CheckForError(err)
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	tlsconfig.InsecureSkipVerify = CRTInsecureSkipVerify
	tlsconfig.ServerName = CRTValidFor
	tlsconfig.RootCAs = pool
	return tlsconfig
}
