package gosyncmodules

import (
	"github.com/go-ldap/ldap"
	"fmt"
	"time"
)

func ConnectToAD(ADHost, AD_Port string, ADUsername string, ADPassword string, ADConnTimeout int) (*ldap.Conn){
	ldap.DefaultTimeout = time.Duration(ADConnTimeout) * time.Second
	Info.Println("Set connection timeout to", ldap.DefaultTimeout)
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", ADHost, AD_Port))
	CheckForError(err)
	Info.Println("Binding")
	err = l.Bind(ADUsername, ADPassword)
	CheckForError(err)
	return l
}
