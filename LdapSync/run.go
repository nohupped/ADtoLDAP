package main


import (
	"fmt"
	"github.com/mavricknz/ldap"
	modules "faimodules"
	"os/user"
	"SyncModules"
	
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
	ldap_port := uint16(ADPort.MustInt(389));

	attributes := make([]string, 0, 1)
	for _, i := range ADAttr.Strings(","){
		attributes = append(attributes, i)
	}
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
	
	//search_request := ldap.NewSearchRequest("ou=Tapestry,dc=internal,dc=media,dc=net", ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, false, "", "", nil)
	
	search_request := ldap.NewSearchRequest(ADBaseDN.String(), ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0,false, "(cn=*)", attributes, nil)
	
	sr, err := connection.SearchWithPaging(search_request, uint32(ADPage.MustInt(100)))
	for _, i := range sr.Entries {
		fmt.Println(i.String())
		fmt.Println("\n\n\n")
	}
	fmt.Println("Result gathered")
	
	
	
	
	
	
	
}