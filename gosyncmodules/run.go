package gosyncmodules

func InitialrunAD(ADHost, AD_Port, ADUsername, ADPassword, ADBaseDN, ADFilter string, ADAttribute []string, ADPage int, ADConnTimeout int, shutdownChannel chan string)  {
	connectAD := ConnectToAD(ADHost, AD_Port, ADUsername, ADPassword, ADConnTimeout)
	defer func() {shutdownChannel <- "Done from func InitialrunAD"}()
	defer Info.Println("closed")
	defer connectAD.Close()
	defer Info.Println("Closing connection")
	ADElements := GetFromAD(connectAD, ADBaseDN, ADFilter, ADAttribute, uint32(ADPage))
	Info.Println(ADElements)

}

func InitialrunLDAP(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPBaseDN, LDAPFilter string, LDAPAttribute []string, LDAPPage int, LDAPConnTimeout int, shutdownChannel chan string)  {
	connectLDAP := ConnectToLdap(LDAPHost, LDAP_Port, LDAPUsername, LDAPPassword, LDAPConnTimeout)
	defer func() {shutdownChannel <- "Done from func InitialrunLDAP"}()
	defer Info.Println("closed")
	defer connectLDAP.Close()
	defer Info.Println("Closing connection")
	Info.Println("End of LDAP")

}
