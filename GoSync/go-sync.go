package main


import (
	"os/user"
	"fmt"
	"os"
	"gosyncmodules"
)

func init()  {
	if len(os.Args) == 1 {
		fmt.Println("Usage:\n\t", os.Args[0],
			"--init to do an init run to freshly populate ldap",
			"from AD. Caution, this can overwrite data. \n\t", os.Args[0],
			"--sync to keep monitoring the changes and sync")
		os.Exit(1)
	}
}

func main() {
	logfileMain := "/var/log/ldapsync.log"
	TAG := gosyncmodules.RandomGen(5)
	username, err := user.Current()
	gosyncmodules.CheckForError(err)
	loggerMain := gosyncmodules.StartLog(logfileMain, username, TAG)
	defer loggerMain.Close()
	configFile := "/etc/ldapsync.ini"
	gosyncmodules.CheckPerm(configFile)
	config, err := gosyncmodules.GetConfig(configFile)
	gosyncmodules.CheckForError(err)

	//AD Variables
	ADGlobal, err := config.GetSection("ADServer")
	gosyncmodules.CheckForError(err)
	ADHost, err := ADGlobal.GetKey("ADHost")
	gosyncmodules.CheckForError(err)
	ADPort, err := ADGlobal.GetKey("ADPort")
	gosyncmodules.CheckForError(err)
	ADPage, err := ADGlobal.GetKey("ADPage")
	gosyncmodules.CheckForError(err)
	ADConnTimeOut, err := ADGlobal.GetKey("ADConnTimeOut")
	gosyncmodules.CheckForError(err)
	ADUsername, err := ADGlobal.GetKey("username")
	gosyncmodules.CheckForError(err)
	ADPassword, err := ADGlobal.GetKey("password")
	gosyncmodules.CheckForError(err)
	ADBaseDN, err := ADGlobal.GetKey("basedn")
	gosyncmodules.CheckForError(err)
	ADAttr, err := ADGlobal.GetKey("attr")
	gosyncmodules.CheckForError(err)
	ADFilter, err := ADGlobal.GetKey("filter")
	gosyncmodules.CheckForError(err)
	AD_Port := ADPort.MustString("389")
	ADAttribute := make([]string, 0, 1)
	for _, i := range ADAttr.Strings(",") {
		ADAttribute = append(ADAttribute, i)
	}
	//LDAP Variables
	LDAPGlobal, err := config.GetSection("LDAPServer")
	gosyncmodules.CheckForError(err)
	LDAPHost, err := LDAPGlobal.GetKey("LDAPHost")
	gosyncmodules.CheckForError(err)
	fmt.Println(LDAPHost)
	LDAPPort, err := LDAPGlobal.GetKey("LDAPPort")
	gosyncmodules.CheckForError(err)
	LDAP_Port := LDAPPort.MustString("389")
	LDAPPage, err := LDAPGlobal.GetKey("LDAPPage")
	gosyncmodules.CheckForError(err)
	LDAPConnTimeOut, err := LDAPGlobal.GetKey("LDAPConnTimeOut")
	gosyncmodules.CheckForError(err)
	LDAPUsername, err := LDAPGlobal.GetKey("username")
	gosyncmodules.CheckForError(err)
	LDAPPassword, err := LDAPGlobal.GetKey("password")
	gosyncmodules.CheckForError(err)
	LDAPBaseDN, err := LDAPGlobal.GetKey("basedn")
	gosyncmodules.CheckForError(err)
	LDAPFilter, err := ADGlobal.GetKey("filter")
	gosyncmodules.CheckForError(err)
	LDAPAttr, err := LDAPGlobal.GetKey("attr")
	gosyncmodules.CheckForError(err)
	LDAPAttribute := make([]string, 0, 1)
	for _, i := range LDAPAttr.Strings(",") {
		LDAPAttribute = append(LDAPAttribute, i)
	}
	//End of variable declaration

	gosyncmodules.Info.Println("\n\tADHost: ", ADHost, "\n\tADPort: ", ADPort, "\n\tADPageSize: ",
		ADPage, "\n\tADBaseDN: ", ADBaseDN, "\n\tADAttr: ", ADAttribute, "\n\tADFilter: ", ADFilter)
	gosyncmodules.Info.Println("\n\tLDAPHost: ", LDAPHost, "\n\tADPort: ", ADPort, "\n\tADPageSize: ",
		ADPage, "\n\tADBaseDN: ", ADBaseDN, "\n\tADAttr: ", ADAttr, "\n\tADFilter: ", ADFilter)
	var howtorun string
	if os.Args[1] == "--init" {
		howtorun = "init"
	} else if os.Args[1] == "--sync" {
		howtorun = "sync"
	} else {
		fmt.Println("Usage:\n\t", os.Args[0],
			"--init to do an init run to freshly populate ldap",
			"from AD. Caution, this can overwrite data. \n\t", os.Args[0],
			"--sync to keep monitoring the changes and sync")
		os.Exit(1)
	}
	gosyncmodules.Info.Println("Starting script with", howtorun, "parameter")

	if howtorun == "init" {
		shutdownChannel := make(chan string, 2)
		defer gosyncmodules.Info.Println()
		defer close(shutdownChannel)
		gosyncmodules.Info.Println("Initializing bool channel and getting AD entries in goroutine")
		gosyncmodules.Info.Println("Gathering results")
		go gosyncmodules.InitialrunAD(ADHost.String(), AD_Port, ADUsername.String(), ADPassword.String(),
			ADBaseDN.String(), ADFilter.String(), ADAttribute, ADPage.MustInt(500), ADConnTimeOut.MustInt(10), shutdownChannel)
		gosyncmodules.InitialrunLDAP(LDAPHost.String(), LDAP_Port, LDAPUsername.String(), LDAPPassword.String(),
			LDAPBaseDN.String(), LDAPFilter.String(), LDAPAttribute, LDAPPage.MustInt(500), LDAPConnTimeOut.MustInt(10), shutdownChannel)
		gosyncmodules.Info.Println(<-shutdownChannel)
		gosyncmodules.Info.Println(<-shutdownChannel)

	}else {
		fmt.Println("No init")
	}
}