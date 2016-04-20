package gosyncmodules

import (
	"github.com/nohupped/ldap" //using a forked version that includes custom methods to retrieve and edit *AddRequest struct.
	"fmt"
	"time"
)

func ConnectToLdap(ADHost, AD_Port string, ADUsername, ADPassword string, ADConnTimeout int) (*ldap.Conn) {
	ldap.DefaultTimeout = time.Duration(ADConnTimeout) * time.Second
	Info.Println("Set LDAP connection timeout to", ldap.DefaultTimeout)
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ADHost, AD_Port))
	CheckForError(err)
	Info.Println("Binding")
	err = l.Bind(ADUsername, ADPassword)
	CheckForError(err)
	return l
}
