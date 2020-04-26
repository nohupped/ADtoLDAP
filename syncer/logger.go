package syncer

import (
	log "github.com/nohupped/glog"
	"os"
)

var logger *log.Logger
var loggerfile *os.File
var logFileOpenErr error

// LoggerClose closes the logger fd.
func LoggerClose() {
	loggerfile.Close()
}

// StartLog starts the logger. Accepts a path to the logfile.
func StartLog(logfile string) *log.Logger {
	loggerfile, logFileOpenErr = os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if logFileOpenErr != nil {
		panic(logFileOpenErr)
	}
	logger = log.New(loggerfile, "", log.Lshortfile)
	l := log.DebugLevel
	logger.SetLogLevel(&l)

	return logger
}

// SetLogLevel sets the loglevel.
func SetLogLevel(loglevel *string) {
	switch *loglevel {
	case "ErrorLevel":
		l := log.ErrorLevel
		logger.SetLogLevel(&l)
		break
	case "WarnLevel":
		l := log.WarnLevel
		logger.SetLogLevel(&l)
		break
	case "InfoLevel":
		l := log.InfoLevel
		logger.SetLogLevel(&l)
		break
	case "DebugLevel":
		l := log.DebugLevel
		logger.SetLogLevel(&l)
		break
	default:
		logger.Warnln("Specified loglevel is invalid. Running with daemon set debug loglevel")
	}
}


