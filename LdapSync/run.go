package main


import (
	modules "faimodules"
	"os/user"
	"SyncModules"
	"fmt"
)

func main() {
	configFile := "/etc/ldapsync.ini"
	logfile := "/var/log/ldapsync.log"
	username, _ := user.Current()
	logger := modules.StartLog(logfile, username)
	defer logger.Close()
	config, _ := SyncModules.GetConfig(configFile)
	ADGlobal, _ := config.GetSection("ADServer")
	ADHost, _ := ADGlobal.GetKey("ADHost")
	ADPort, _ := ADGlobal.GetKey("ADPort")
	ADPage, _ := ADGlobal.GetKey("ADPage")
	ADUsername, _ := ADGlobal.GetKey("username")
	ADPassword, _ := ADGlobal.GetKey("password")
	ADBaseDN, _ := ADGlobal.GetKey("basedn")
	ADAttr, _ := ADGlobal.GetKey("attr")
	ADFilter, _ := ADGlobal.GetKey("filter")
	ldap_port := uint16(ADPort.MustInt(389))


	ADAttribute := make([]string, 0, 1)
	for _, i := range ADAttr.Strings(","){
		ADAttribute = append(ADAttribute, i)
	}
<<<<<<< HEAD
	modules.Info.Println("\n\tADHost: ", ADHost, "\n\tADPort: ", ADPort, "\n\tADPageSize: ",
		ADPage, "\n\tADBaseDN: ", ADBaseDN, "\n\tADAttr: ", ADAttr, "\n\tADFilter: ", ADFilter)

	connect := SyncModules.ConnectToAD(ADHost.String(), ldap_port, ADUsername.String(), ADPassword.String())
	defer modules.Info.Println(connect.Addr, "closed")
	defer connect.Close()
	defer modules.Info.Println("Closing connection", connect.Addr)
	ADElements := SyncModules.GetFromAD(connect, ADBaseDN.String(), ADFilter.String(), ADAttribute, uint32(ADPage.MustInt(500)))
	for _, x := range ADElements {
		fmt.Println(x.DN)
=======
	modules.Info.Println(attributes)
	modules.Info.Println("Connecting to ldap server", ADHost.String())
	connection := ldap.NewLDAPConnection(ADHost.String(), ldap_port)
	modules.Info.Println(connection)
	//Connect
	err := connection.Connect()
	if err != nil{
		modules.Error.Println(err)
	}
	defer connection.Close()
	//Bind
	err = connection.Bind(ADUsername.String(), ADPassword.String())
	if err != nil {
		modules.Error.Println(err)
	}
	modules.Info.Println("Binded")
	
	search_request := ldap.NewSearchRequest(ADBaseDN.String(), ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0,false, "(cn=*)", attributes, nil)
	
	sr, err := connection.SearchWithPaging(search_request, uint32(ADPage.MustInt(100)))
	for _, i := range sr.Entries {
		fmt.Println(i.String())
		fmt.Println("\n\n\n")
>>>>>>> 404c6ff819a5a4cfbaf9bdc47fda2ba472e5674a
	}
	
}
