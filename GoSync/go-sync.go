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
	username, _ := user.Current()
	loggerMain := gosyncmodules.StartLog(logfileMain, username, TAG)
	defer loggerMain.Close()
	configFile := "/etc/ldapsync.ini"
	config, _ := gosyncmodules.GetConfig(configFile)
	ADGlobal, _ := config.GetSection("ADServer")
	ADHost, _ := ADGlobal.GetKey("ADHost")
	ADPort, _ := ADGlobal.GetKey("ADPort")
	ADPage, _ := ADGlobal.GetKey("ADPage")
	ADConnTimeOut, _ := ADGlobal.GetKey("ADConnTimeOut")
	ADUsername, _ := ADGlobal.GetKey("username")
	ADPassword, _ := ADGlobal.GetKey("password")
	ADBaseDN, _ := ADGlobal.GetKey("basedn")
	ADAttr, _ := ADGlobal.GetKey("attr")
	ADFilter, _ := ADGlobal.GetKey("filter")
	AD_Port := ADPort.MustString("389")
	ADAttribute := make([]string, 0, 1)
	for _, i := range ADAttr.Strings(",") {
		ADAttribute = append(ADAttribute, i)
	}
	gosyncmodules.Info.Println("\n\tADHost: ", ADHost, "\n\tADPort: ", ADPort, "\n\tADPageSize: ",
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
		shutdownChannel := make(chan bool)
		defer close(shutdownChannel)
		gosyncmodules.Info.Println("Initializing bool channel and getting AD entries in goroutine")
		gosyncmodules.Info.Println("Gathering results")
		go gosyncmodules.Initialrun(ADHost.String(), AD_Port, ADUsername.String(), ADPassword.String(),
			ADBaseDN.String(), ADFilter.String(), ADAttribute, ADPage.MustInt(500), ADConnTimeOut.MustInt(10), shutdownChannel)
		gosyncmodules.Info.Println(<- shutdownChannel)

	}else {
		fmt.Println("No init")
	}

}