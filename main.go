package main

import (
	"os/user"
	"github.com/nohupped/ADtoLDAP/syncer"
	"os"
	"reflect"
	"flag"
	"gopkg.in/ldap.v2"
	"runtime"
	"time"
	"fmt"
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

	checkSafety := flag.Bool("safe", true, "Set it to false to skip config file securitycheck")
	configFile := flag.String("configfile", "/etc/ldapsync.ini", "Path to the config file")
	logfile := flag.String("logfile", "/var/log/ldapsync.log", "Path to log file. Defaults to /var/log/ldapsync.log")
	sampleConfig := flag.Bool("showSampleConfig", false, "Prints a sample config to STDOUT")
	flag.Parse()
	if *sampleConfig == true {
		fmt.Println(syncer.SampleConfig)
		os.Exit(0)
	}

	username, err := user.Current()
	syncer.CheckForError(err)
	log := syncer.StartLog(*logfile)
	defer syncer.LoggerClose()
	log.Infoln("Running program as", username)

	log.Infoln("safe option set to", *checkSafety)
	log.Infoln("Config file is, ", *configFile)

	if *checkSafety == true {
		syncer.CheckPerm(*configFile)
	} else {
		log.Infoln("Skipping file permission check on", *configFile)
	}
	r := syncer.NewRuntimeConfig(*configFile)

	//End of variable declaration

	log.Debugf("%+v\n", *r)

	log.Infoln("Initiating sync")

	for {
		AddChan := make(chan syncer.Action)
		log.Debugln("Created", reflect.TypeOf(AddChan))
		DelChan := make(chan syncer.Action)
		log.Debugln("Created", reflect.TypeOf(DelChan))
		shutdownAddChan := make(chan string)
		shutdownDelChan := make(chan string)
		shutdownChannel := make(chan string)
		LdapConnectionChan := make(chan *ldap.Conn)
		log.Debugln("Created channel of type", reflect.TypeOf(LdapConnectionChan))
		ADElementsChan := make(chan *[]syncer.LDAPElement)
		log.Debugln("Created channel of type", reflect.TypeOf(ADElementsChan))
		LDAPElementsChan := make(chan *[]syncer.LDAPElement)
		log.Debugln("Created channel of type", reflect.TypeOf(LDAPElementsChan))

		go syncer.SyncRunLDAP(r.LDAPServer.Host, r.LDAPServer.Port, r.LDAPServer.Username, r.LDAPServer.Password,
			r.LDAPServer.BaseDN, r.LDAPServer.Filter, r.LDAPServer.Attributes, r.LDAPServer.Page,
			r.LDAPServer.ConnTimeOut, r.LDAPServer.UseTLS, r.LDAPServer.CRTInsecureSkipVerify,
			r.LDAPServer.CRTValidFor, r.LDAPServer.CRTPath, shutdownChannel, LDAPElementsChan, LdapConnectionChan,
			r.ReplaceAttributes, r.MapAttributes)
		go syncer.SyncRunAD(r.ADServer.Host, r.ADServer.Port, r.ADServer.Username, r.ADServer.Password,
			r.ADServer.BaseDN, r.ADServer.Filter, r.ADServer.Attributes, r.ADServer.Page,
			r.ADServer.ConnTimeOut, r.ADServer.UseTLS, r.ADServer.CRTInsecureSkipVerify,
			r.ADServer.CRTValidFor, r.ADServer.CRTPath, shutdownChannel, ADElementsChan)
		ADElements := <-ADElementsChan
		LDAPElements := <-LDAPElementsChan
		LDAPConnection := <-LdapConnectionChan
		log.Debugln(<-shutdownChannel)
		log.Debugln(<-shutdownChannel)

		ADElementsConverted := syncer.InitialPopulateToLdap(ADElements, LDAPConnection, r.ReplaceAttributes, r.MapAttributes, true)
		LDAPElementsConverted := syncer.InitialPopulateToLdap(LDAPElements, LDAPConnection, r.ReplaceAttributes, r.MapAttributes, true)

		syncer.ConvertRealmToLower(ADElementsConverted)
		log.Debugln("Converted AD Realms to lowercase")

		go syncer.FindAdds(&ADElementsConverted, &LDAPElementsConverted, LDAPConnection, AddChan, shutdownAddChan)
		go syncer.FindDels(&LDAPElementsConverted, &ADElementsConverted, DelChan, shutdownDelChan)
		counter := 0
		for {
			select {
			case Add := <-AddChan:
				for k, v := range Add {
					log.Debugln(k, ":", v)
					err := LDAPConnection.Add(v)
					if err != nil {
						log.Errorln(err)
					}
				}
			case Del := <-DelChan:
				for k, v := range Del {
					log.Debugln(k, ":", v)
					deleteRecord := ldap.NewDelRequest(v.DN, []ldap.Control{})
					err := LDAPConnection.Del(deleteRecord)
					if err != nil {
						log.Errorln(err)
					}
				}
			case shutdownAdd := <-shutdownAddChan:
				log.Debugln(shutdownAdd)
				counter ++
			case shutdownDel := <-shutdownDelChan:
				log.Debugln(shutdownDel)
				counter ++

			}
			if counter == 2 {
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
		log.Infoln("Sleeping for", r.Delay, "seconds, and iterating again...")
		log.Infoln("Currently active goroutines: ", runtime.NumGoroutine())
		//Thanks to profiling, that helped finding a goroutine leak.
		/*buf1 := new(bytes.Buffer)
		pprof.Lookup("goroutine").WriteTo(buf1, 1)
		fmt.Println("pprof.Lookup.WriteTo report:\n", string(buf1.Bytes()[:buf1.Len()]))
		var buf [10240]byte
		out := buf[:runtime.Stack(buf[:], true)]
		fmt.Println("runtime.Stack report:\n", string(out))*/
		time.Sleep(time.Second * time.Duration(r.Delay))

	}
}
