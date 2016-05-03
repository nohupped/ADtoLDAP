package main


import (
	"os/user"
	"fmt"
	"os"
//	"gosyncmodules"
	"github.com/nohupped/ADtoLDAP/gosyncmodules"
	"reflect"
//	"os/signal"
	"github.com/nohupped/ldap" //using a forked version that includes custom methods to retrieve and edit *AddRequest struct.
	"time"
	"runtime"
//	"bytes"
//	"runtime/pprof"
)

func init()  {
	if len(os.Args) == 1 {
		fmt.Println("Usage:\n\t", os.Args[0],
			"--init to do an init run to freshly populate ldap",
			"from AD. Caution, this can overwrite data. \n\t", os.Args[0],
			"--sync to keep monitoring the changes and sync")
		os.Exit(1)
	}

/*	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			fmt.Println(sig.String(), "received, terminating.")
			os.Exit(1)
		}
	}()*/
}

func main() {
	logfileMain := "/var/log/ldapsync.log"
	//TAG := gosyncmodules.RandomGen(5)
	username, err := user.Current()
	gosyncmodules.CheckForError(err)
	loggerMain := gosyncmodules.StartLog(logfileMain, username)
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

	//AD to LDAP Mapping, replace and required variables
	ReplaceAttributes, err := config.GetSection("Replace")
	gosyncmodules.CheckForError(err)
	MapAttributes, err := config.GetSection("Map")
	gosyncmodules.CheckForError(err)

	//SyncVariables

	SyncDelay, err := config.GetSection("Sync")
	gosyncmodules.CheckForError(err)
	Delay, err := SyncDelay.GetKey("sleepTime")
	gosyncmodules.CheckForError(err)
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
			"from AD. Caution, if there is data already populated in ldap,\n\t\t",
			"you may have to wipe it clean before doing init. \n\t", os.Args[0],
			"--sync to keep monitoring the changes and sync")
		os.Exit(1)
	}
	gosyncmodules.Info.Println("Starting script with", howtorun, "parameter")

	if howtorun == "init" {
		shutdownChannel := make(chan string)
		defer gosyncmodules.Info.Println("Closed blocking channel")
		defer close(shutdownChannel)
		gosyncmodules.Info.Println("Initializing bool channel and getting AD entries in goroutine")
		gosyncmodules.Info.Println("Gathering results")

		//Create channel to receive slice of struct
		ADElementsChan := make(chan *[]gosyncmodules.LDAPElement)
		gosyncmodules.Info.Println("Created channel of type", reflect.TypeOf(ADElementsChan))

		go gosyncmodules.InitialrunAD(ADHost.String(), AD_Port, ADUsername.String(), ADPassword.String(),
			ADBaseDN.String(), ADFilter.String(), ADAttribute, ADPage.MustInt(500), ADConnTimeOut.MustInt(10), shutdownChannel, ADElementsChan)
		ADElements := <- ADElementsChan		//Finished retriving AD results
		gosyncmodules.Info.Println(<-shutdownChannel)	//Finished reading from Blocking channel

		gosyncmodules.InitialrunLDAP(LDAPHost.String(), LDAP_Port, LDAPUsername.String(), LDAPPassword.String(),
			LDAPBaseDN.String(), LDAPFilter.String(), LDAPAttribute, LDAPPage.MustInt(500), LDAPConnTimeOut.MustInt(10), ADElements,
			ReplaceAttributes, MapAttributes)

		gosyncmodules.Info.Println("Received", reflect.TypeOf(ADElementsChan), "from child thread, and has ", len(*ADElements), "elements")

	}else {
		gosyncmodules.Info.Println("Initiating sync")

		for ; ;  {
			AddChan := make(chan gosyncmodules.Action)
			gosyncmodules.Info.Println("Created", reflect.TypeOf(AddChan))
			DelChan := make(chan gosyncmodules.Action)
			gosyncmodules.Info.Println("Created", reflect.TypeOf(DelChan))
			shutdownAddChan := make(chan string)
			shutdownDelChan := make(chan string)
			shutdownChannel := make(chan string)
			LdapConnectionChan := make(chan *ldap.Conn)
			gosyncmodules.Info.Println("Created channel of type", reflect.TypeOf(LdapConnectionChan))
			ADElementsChan := make(chan *[]gosyncmodules.LDAPElement)
			gosyncmodules.Info.Println("Created channel of type", reflect.TypeOf(ADElementsChan))
			LDAPElementsChan := make(chan *[]gosyncmodules.LDAPElement)
			gosyncmodules.Info.Println("Created channel of type", reflect.TypeOf(LDAPElementsChan))


			go gosyncmodules.SyncrunLDAP(LDAPHost.String(), LDAP_Port, LDAPUsername.String(), LDAPPassword.String(),
					LDAPBaseDN.String(), LDAPFilter.String(), LDAPAttribute, LDAPPage.MustInt(500),
					LDAPConnTimeOut.MustInt(10), shutdownChannel, LDAPElementsChan, LdapConnectionChan,
					ReplaceAttributes, MapAttributes)
			go gosyncmodules.InitialrunAD(ADHost.String(), AD_Port, ADUsername.String(), ADPassword.String(),
				ADBaseDN.String(), ADFilter.String(), ADAttribute, ADPage.MustInt(500),
				ADConnTimeOut.MustInt(10), shutdownChannel, ADElementsChan)
			ADElements := <- ADElementsChan
			LDAPElements := <- LDAPElementsChan
			LDAPConnection := <- LdapConnectionChan
			gosyncmodules.Info.Println(<-shutdownChannel)
			gosyncmodules.Info.Println(<-shutdownChannel)

			ADElementsConverted := gosyncmodules.InitialPopulateToLdap(ADElements, LDAPConnection, ReplaceAttributes, MapAttributes, true)
			LDAPElementsConverted := gosyncmodules.InitialPopulateToLdap(LDAPElements, LDAPConnection, ReplaceAttributes, MapAttributes, true)



			gosyncmodules.ConvertRealmToLower(ADElementsConverted)
			gosyncmodules.Info.Println("Converted AD Realms to lowercase")



			go gosyncmodules.FindAdds(&ADElementsConverted, &LDAPElementsConverted, LDAPConnection, AddChan, shutdownAddChan)
			go gosyncmodules.FindDels(&LDAPElementsConverted, &ADElementsConverted, DelChan, shutdownDelChan)
			counter := 0
			for ; ; {
				select {
				case Add := <- AddChan:
					for _, v := range Add {
					//	gosyncmodules.Info.Println(k, ":", v)
						err := LDAPConnection.Add(v)
						if err != nil {
							//gosyncmodules.Error.Println(err)
						}
					}
				case Del := <- DelChan:
					for _, v := range Del  {
					//	gosyncmodules.Info.Println(k, ":", v)
						delete := ldap.NewDelRequest(v.DN, []ldap.Control{})
						err := LDAPConnection.Del(delete)
						if err != nil {
							//gosyncmodules.Error.Println(err)
						}
					}
				case shutdownAdd := <- shutdownAddChan:
					gosyncmodules.Info.Println(shutdownAdd)
					counter += 1
				case shutdownDel := <- shutdownDelChan:
					gosyncmodules.Info.Println(shutdownDel)
					counter += 1

				}
				if counter == 2{
					gosyncmodules.Info.Println("Counter reached")
					break
				}

			}
			LDAPConnection.Close()
			//Sleep the daemon
			close(shutdownDelChan)
			close(shutdownAddChan)
			close(shutdownChannel)
			close(LdapConnectionChan)
			close(ADElementsChan)
			close(LDAPElementsChan)
			gosyncmodules.Info.Println("Sleeping for", Delay.MustInt(5), "seconds, and iterating again...")
			gosyncmodules.Info.Println("Current active goroutines: ", runtime.NumGoroutine())
			//Thanks to profiling, that helped finding a goroutine leak.
			/*buf1 := new(bytes.Buffer)
			pprof.Lookup("goroutine").WriteTo(buf1, 1)
			fmt.Println("pprof.Lookup.WriteTo report:\n", string(buf1.Bytes()[:buf1.Len()]))
			var buf [10240]byte
			out := buf[:runtime.Stack(buf[:], true)]
			fmt.Println("runtime.Stack report:\n", string(out))*/
			time.Sleep(time.Second * time.Duration(Delay.MustInt(5)))



		}
	}
}
