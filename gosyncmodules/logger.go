package gosyncmodules

import (
	"fmt"
	"log"
	"io"
	"os"
	"os/user"
	"strings"
	"math/rand"
	"time"
)


var (
	Trace *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

func RandomGen(length int) string  {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func logInit(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer, user *user.User, TAG string, tmp []string) {

	Trace = log.New(traceHandle,
		"TRACE: "+TAG+"  "+user.Username+"  " + strings.Join(tmp, " ") + " ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: "+TAG+"  "+user.Username+"  " + strings.Join(tmp, " ") + " ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: "+TAG+"  "+user.Username+"  " + strings.Join(tmp, " ") + " ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: "+TAG+"  "+user.Username+"  " + strings.Join(tmp, " ") + " ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

func StartLog(logfile string, user *user.User, TAG string,  tmp ...string) *os.File{
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Cannot open logfile")
		panic(err)
	}
	multi := io.MultiWriter(file, os.Stdout)

	logInit(multi, multi, multi, multi, user, TAG, tmp)
	Trace.SetOutput(file)
	Info.SetOutput(file)
	Warning.SetOutput(file)
	Error.SetOutput(file)
	return file
}

