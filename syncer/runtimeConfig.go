package syncer

import (
	"gopkg.in/ini.v1"
)

// getConfig is used to read the config from file and returns an *ini.File object and error.
func getConfig(configFile string) (*ini.File, error) {
	var Cfg *ini.File
	Cfg, err := ini.Load(configFile)
	if err != nil {
		return Cfg, err
	}
	return Cfg, err
}

// RuntimeConfig stores the parsed ini config from file as a struct.
type RuntimeConfig struct {
	//	From Server section
	ADServer                         DS
	LDAPServer                       DS
	ReplaceAttributes, MapAttributes *ini.Section
	Delay                            int
	LogLevel                         string
}

// DS is used to store the directory server config
type DS struct {
	Host, Port, CRTValidFor, CRTPath, Username, Password,
	BaseDN, Filter string
	UseTLS, CRTInsecureSkipVerify bool
	Page, ConnTimeOut             int
	Attributes                    []string
}

// NewRuntimeConfig is used to return a RuntimeConfig struct. This is populated by reading the server's
// configuration ini file (default: /etc/ldapsync.ini) and parsing it.
func NewRuntimeConfig(path string) *RuntimeConfig {

	r := new(RuntimeConfig)

	config, err := getConfig(path)
	CheckForError(err)

	getDSSections("ADServer", config, &r.ADServer)
	getDSSections("LDAPServer", config, &r.LDAPServer)

	//AD to LDAP Mapping, replace and required variables
	r.ReplaceAttributes, err = config.GetSection("Replace")
	CheckForError(err)
	r.MapAttributes, err = config.GetSection("Map")
	CheckForError(err)

	// Daemon settings
	DaemonConfig, err := config.GetSection("Sync")
	CheckForError(err)
	Delay, err := DaemonConfig.GetKey("sleepTime")
	CheckForError(err)
	r.Delay = Delay.MustInt(5)

	if config.Section("Sync").HasKey("loglevel") {
		level, err := DaemonConfig.GetKey("loglevel")
		var loglevel *string
		if err != nil {
			logger.Warnln(err, "encountered, using system-set debug loglevel")
		} else {
			logLevelP := level.String()
			loglevel = &logLevelP
		}
		SetLogLevel(loglevel)
	} else {
		logger.Warnln("Loglevel not defined in config", path, ". Using system-set DEBUG level.")
	}

	return r
}

func getDSSections(section string, config *ini.File, r *DS) {
	Global, err := config.GetSection(section)
	CheckForError(err)

	Host, err := Global.GetKey("Host")
	CheckForError(err)
	r.Host = Host.String()

	Port, err := Global.GetKey("Port")
	CheckForError(err)
	r.Port = Port.MustString("389")

	UseTLS, err := Global.GetKey("UseTLS")
	CheckForError(err)
	r.UseTLS = UseTLS.MustBool(true)

	switch r.UseTLS {
	case true:

		CrtValidFor, err := Global.GetKey("CRTValidFor")
		CheckForError(err)
		r.CRTValidFor = CrtValidFor.String()

		CrtPath, err := Global.GetKey("CRTPath")
		CheckForError(err)
		r.CRTPath = CrtPath.String()

		CRTInsecureSkipVerify, err := Global.GetKey("InsecureSkipVerify")
		CheckForError(err)
		r.CRTInsecureSkipVerify = CRTInsecureSkipVerify.MustBool(false)

	case false:

		r.CRTValidFor = "DUMMY"
		r.CRTPath = "DUMMY"
		r.CRTInsecureSkipVerify = false

	}

	Page, err := Global.GetKey("Page")
	CheckForError(err)
	r.Page = Page.MustInt(500)

	ConnTimeOut, err := Global.GetKey("ConnTimeOut")
	CheckForError(err)
	r.ConnTimeOut = ConnTimeOut.MustInt(10)

	Username, err := Global.GetKey("username")
	CheckForError(err)
	r.Username = Username.String()

	Password, err := Global.GetKey("password")
	CheckForError(err)
	r.Password = Password.String()

	BaseDN, err := Global.GetKey("basedn")
	CheckForError(err)
	r.BaseDN = BaseDN.String()

	Attr, err := Global.GetKey("attr")
	CheckForError(err)
	//Attributes := make([]string, 0, 1)
	for _, i := range Attr.Strings(",") {
		r.Attributes = append(r.Attributes, i)
	}

	Filter, err := Global.GetKey("filter")
	CheckForError(err)
	r.Filter = Filter.String()

}
