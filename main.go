package main


import (
	"os/user"
//	"fmt"
	"os"
	"ADtoLDAP/gosyncmodules"
//	"github.com/nohupped/ADtoLDAP/gosyncmodules"
	"reflect"
//	"os/signal"
	"gopkg.in/ldap.v2"
	"time"
	"runtime"
	"flag"
//	"bytes"
//	"runtime/pprof"
)

/*func init()  {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			fmt.Println(sig.String(), "received, terminating.")
			os.Exit(1)
		}
	}()
}*/

func main() {

	// Initialize logger


	// Flags
	checkSafety := flag.Bool("safe", true, "Set it to false to skip config file securitycheck")
	syncrun := flag.String("sync", "daemon", "Set it to \"once\" for a single run, and \"daemon\" to run it continuously")
	configFile := flag.String("configfile", "/etc/ldapsync.ini", "Path to the config file")
	logfile := flag.String("logfile", "/var/log/ldapsync.log", "Path to log file. Defaults to /var/log/ldapsync.log")
	flag.Parse()

	username, err := user.Current()
	gosyncmodules.CheckForError(err)
	log := gosyncmodules.StartLog(*logfile)
	defer gosyncmodules.LoggerClose()
	log.Infoln("Running program as", username)

	log.Infoln("safe option set to", *checkSafety)
	log.Infoln("Config file is, ", *configFile)


	if *checkSafety == true {
		gosyncmodules.CheckPerm(*configFile)
	} else {
		log.Infoln("Skipping file permission check on", *configFile)
	}
	config, err := gosyncmodules.GetConfig(*configFile)
	gosyncmodules.CheckForError(err)

	//AD Variables
	ADGlobal, err := config.GetSection("ADServer")
	gosyncmodules.CheckForError(err)
	ADHost, err := ADGlobal.GetKey("ADHost")
	gosyncmodules.CheckForError(err)
	ADPort, err := ADGlobal.GetKey("ADPort")
	gosyncmodules.CheckForError(err)
	ADUseTLS, err := ADGlobal.GetKey("UseTLS")
	gosyncmodules.CheckForError(err)
	ADCrtValidFor, err := ADGlobal.GetKey("CRTValidFor")
	gosyncmodules.CheckForError(err)
	ADCrtPath, err := ADGlobal.GetKey("CRTPath")
	gosyncmodules.CheckForError(err)
	ADCRTInsecureSkipVerify, err := ADGlobal.GetKey("InsecureSkipVerify")
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
	LDAPUseTLS, err := LDAPGlobal.GetKey("UseTLS")
	gosyncmodules.CheckForError(err)
	LDAPCrtValidFor, err := LDAPGlobal.GetKey("CRTValidFor")
	gosyncmodules.CheckForError(err)
	LDAPCrtPath, err := LDAPGlobal.GetKey("CRTPath")
	gosyncmodules.CheckForError(err)
	LDAPCRTInsecureSkipVerify, err := LDAPGlobal.GetKey("InsecureSkipVerify")
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

	Daemon_config, err := config.GetSection("Sync")
	gosyncmodules.CheckForError(err)
	Delay, err := Daemon_config.GetKey("sleepTime")
	gosyncmodules.CheckForError(err)
	if config.Section("Sync").HasKey("loglevel") {
		l_level, err := Daemon_config.GetKey("loglevel")
		var loglevel *string
		if err != nil {
			log.Warnln(err, "encountered, using system-set debug loglevel")
		} else {
			log_level_p := l_level.String()
			loglevel = &log_level_p
		}
		gosyncmodules.SetLogLevel(loglevel)
	} else {
		log.Warnln("Loglevel not defined in config", *configFile, "so using system-set DEBUG level.")
	}



	//End of variable declaration
	log.Infoln("ADHost: ", ADHost)
	log.Infoln("ADPort: ", ADPort)
	log.Infoln("ADPageSize: ", ADPage)
	log.Infoln("ADBaseDN: ", ADBaseDN)
	log.Infoln("ADAttr: ", ADAttribute)
	log.Infoln("ADFilter: ", ADFilter)

	log.Infoln("LDAPHost: ", LDAPHost)
	log.Infoln("LDAPPort: ", LDAP_Port)
	log.Infoln("LDAPPageSize: ", LDAPPage)
	log.Infoln("LDAPBaseDN: ", LDAPBaseDN)
	log.Infoln("LDAPAttr: ", LDAPAttribute)
	log.Infoln("LDAPFilter: ", LDAPFilter)
	var howtorun string
	if *syncrun == "once" {
		howtorun = "init"
	} else if *syncrun == "daemon" {
		howtorun = "sync"
	} else {

		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Infoln("Starting script with", *syncrun, "parameter")

	if howtorun == "init" {
		shutdownChannel := make(chan string)
		defer log.Infoln("Closed blocking channel")
		defer close(shutdownChannel)
		log.Infoln("Initializing bool channel and getting AD entries in goroutine")
		log.Infoln("Gathering results")

		//Create channel to receive slice of struct
		ADElementsChan := make(chan *[]gosyncmodules.LDAPElement)
		log.Infoln("Created channel of type", reflect.TypeOf(ADElementsChan))

		go gosyncmodules.InitialrunAD(ADHost.String(), AD_Port, ADUsername.String(), ADPassword.String(),
			ADBaseDN.String(), ADFilter.String(), ADAttribute, ADPage.MustInt(500), ADConnTimeOut.MustInt(10),
			ADUseTLS.MustBool(true), ADCRTInsecureSkipVerify.MustBool(false),
			ADCrtValidFor.String(), ADCrtPath.String(), shutdownChannel, ADElementsChan)
		ADElements := <- ADElementsChan		//Finished retriving AD results
		log.Infoln(<-shutdownChannel)	//Finished reading from Blocking channel

		gosyncmodules.InitialrunLDAP(LDAPHost.String(), LDAP_Port, LDAPUsername.String(), LDAPPassword.String(),
			LDAPBaseDN.String(), LDAPFilter.String(), LDAPAttribute, LDAPPage.MustInt(500), LDAPConnTimeOut.MustInt(10),
			LDAPUseTLS.MustBool(true), LDAPCrtValidFor.String(), LDAPCrtPath.String(), LDAPCRTInsecureSkipVerify.MustBool(false),
			ADElements, ReplaceAttributes, MapAttributes)

		log.Infoln("Received", reflect.TypeOf(ADElementsChan), "from child thread, and has ", len(*ADElements), "elements")

	}else {
		log.Infoln("Initiating sync")

		for ; ;  {
			AddChan := make(chan gosyncmodules.Action)
			log.Debugln("Created", reflect.TypeOf(AddChan))
			DelChan := make(chan gosyncmodules.Action)
			log.Debugln("Created", reflect.TypeOf(DelChan))
			shutdownAddChan := make(chan string)
			shutdownDelChan := make(chan string)
			shutdownChannel := make(chan string)
			LdapConnectionChan := make(chan *ldap.Conn)
			log.Debugln("Created channel of type", reflect.TypeOf(LdapConnectionChan))
			ADElementsChan := make(chan *[]gosyncmodules.LDAPElement)
			log.Debugln("Created channel of type", reflect.TypeOf(ADElementsChan))
			LDAPElementsChan := make(chan *[]gosyncmodules.LDAPElement)
			log.Debugln("Created channel of type", reflect.TypeOf(LDAPElementsChan))


			go gosyncmodules.SyncrunLDAP(LDAPHost.String(), LDAP_Port, LDAPUsername.String(), LDAPPassword.String(),
					LDAPBaseDN.String(), LDAPFilter.String(), LDAPAttribute, LDAPPage.MustInt(500),
					LDAPConnTimeOut.MustInt(10), LDAPUseTLS.MustBool(true), LDAPCRTInsecureSkipVerify.MustBool(false),
					LDAPCrtValidFor.String(), LDAPCrtPath.String(), shutdownChannel, LDAPElementsChan, LdapConnectionChan,
					ReplaceAttributes, MapAttributes)
			go gosyncmodules.InitialrunAD(ADHost.String(), AD_Port, ADUsername.String(), ADPassword.String(),
				ADBaseDN.String(), ADFilter.String(), ADAttribute, ADPage.MustInt(500),
				ADConnTimeOut.MustInt(10), ADUseTLS.MustBool(true), ADCRTInsecureSkipVerify.MustBool(false),
				ADCrtValidFor.String(), ADCrtPath.String(), shutdownChannel, ADElementsChan)
			ADElements := <- ADElementsChan
			LDAPElements := <- LDAPElementsChan
			LDAPConnection := <- LdapConnectionChan
			log.Debugln(<-shutdownChannel)
			log.Debugln(<-shutdownChannel)

			ADElementsConverted := gosyncmodules.InitialPopulateToLdap(ADElements, LDAPConnection, ReplaceAttributes, MapAttributes, true)
			LDAPElementsConverted := gosyncmodules.InitialPopulateToLdap(LDAPElements, LDAPConnection, ReplaceAttributes, MapAttributes, true)



			gosyncmodules.ConvertRealmToLower(ADElementsConverted)
			log.Debugln("Converted AD Realms to lowercase")



			go gosyncmodules.FindAdds(&ADElementsConverted, &LDAPElementsConverted, LDAPConnection, AddChan, shutdownAddChan)
			go gosyncmodules.FindDels(&LDAPElementsConverted, &ADElementsConverted, DelChan, shutdownDelChan)
			counter := 0
			for ; ; {
				select {
				case Add := <- AddChan:
					for k, v := range Add {
						log.Debugln(k, ":", v)
						err := LDAPConnection.Add(v)
						if err != nil {
							log.Errorln(err)
						}
					}
				case Del := <- DelChan:
					for k, v := range Del  {
						log.Debugln(k, ":", v)
						delete := ldap.NewDelRequest(v.DN, []ldap.Control{})
						err := LDAPConnection.Del(delete)
						if err != nil {
							log.Errorln(err)
						}
					}
				case shutdownAdd := <- shutdownAddChan:
					log.Debugln(shutdownAdd)
					counter += 1
				case shutdownDel := <- shutdownDelChan:
					log.Debugln(shutdownDel)
					counter += 1

				}
				if counter == 2{
					log.Debugln("Counter reached")
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
			log.Infoln("Sleeping for", Delay.MustInt(5), "seconds, and iterating again...")
			log.Infoln("Currently active goroutines: ", runtime.NumGoroutine())
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
