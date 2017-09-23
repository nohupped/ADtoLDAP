package gosyncmodules

import (
	//"fmt"
	////"log"
	//"io"
	//"os"
	//"os/user"
	//"strings"
	//"math/rand"
	//"time"
	log "github.com/nohupped/glog"
	"os"
)

var logger *log.Logger
var loggerfile *os.File
var logFileOpenErr error

func LoggerClose() {
	loggerfile.Close()
}
func StartLog(logfile string) *log.Logger {
	loggerfile, logFileOpenErr = os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if logFileOpenErr != nil {
		panic(logFileOpenErr)
	}
	logger = log.New(loggerfile, "", log.Lshortfile)
	logger.SetLogLevel(log.DebugLevel)

	return logger
}
func SetLogLevel(loglevel *string) {
	switch *loglevel {
	case "ErrorLevel":
		logger.SetLogLevel(log.ErrorLevel)
		break
	case "WarnLevel":
		logger.SetLogLevel(log.WarnLevel)
		break
	case "InfoLevel":
		logger.SetLogLevel(log.InfoLevel)
		break
	case "DebugLevel":
		logger.SetLogLevel(log.DebugLevel)
		break
	default:
		logger.Warnln("Specified loglevel is invalid. Running with daemon set debug loglevel")
	}
}

//
//var (
//	Trace *log.Logger
//	Info *log.Logger
//	Warning *log.Logger
//	Error *log.Logger
//)
//
//func RandomGen(length int) string  {
//	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
//	b := make([]rune, length)
//	rand.Seed(time.Now().UTC().UnixNano())
//	for i := range b {
//		b[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(b)
//}
//
//
//func logInit(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer, user *user.User, TAG []string) {
//
//	Trace = log.New(traceHandle,
//		"TRACE: "+"  "+user.Username+"  " + strings.Join(TAG, " ") + " ",
//		log.Ldate|log.Ltime|log.Lshortfile)
//
//	Info = log.New(infoHandle,
//		"INFO: "+"  "+user.Username+"  " + strings.Join(TAG, " ") + " ",
//		log.Ldate|log.Ltime|log.Lshortfile)
//
//	Warning = log.New(warningHandle,
//		"WARNING: "+"  "+user.Username+"  " + strings.Join(TAG, " ") + " ",
//		log.Ldate|log.Ltime|log.Lshortfile)
//
//	Error = log.New(errorHandle,
//		"ERROR: "+"  "+user.Username+"  " + strings.Join(TAG, " ") + " ",
//		log.Ldate|log.Ltime|log.Lshortfile)
//
//}
//
//func StartLog(logfile string, user *user.User, TAG ...string) *os.File{
//	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
//	if err != nil {
//		fmt.Println("Cannot open logfile")
//		panic(err)
//	}
//	multi := io.MultiWriter(file, os.Stdout)
//
//	logInit(multi, multi, multi, multi, user, TAG)
//	Trace.SetOutput(file)
//	Info.SetOutput(file)
//	Warning.SetOutput(file)
//	Error.SetOutput(file)
//	return file
//}
//
