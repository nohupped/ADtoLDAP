package gosyncmodules

import (
	"fmt"
	"unsafe"
	"os"
	"os/user"
	"strconv"
)
//#include <sys/stat.h>
//#include <stdlib.h>
import "C"

func CheckPerm(filename string) {
	current_uid := os.Getuid()
	current_gid := os.Getgid()
	username, err := user.LookupId(strconv.Itoa(current_uid))
	var current_username string
	var current_groupname string
	if err != nil{
		current_username= "<Couldn't lookup username>"
	}
	groupname, err := user.LookupGroupId(strconv.Itoa(current_gid))
	if err != nil {
		current_groupname = "<Couldn't lookup groupname>"

	}
	current_username = username.Username
	current_groupname = groupname.Name

	logger.Infoln("using cgo to perform security check on ", filename)
	statstruct := C.stat //stat struct from C
	logger.Infoln("Initiated stat struct")
	path := C.CString(filename)
	logger.Infoln("Converted native string to C.CString")
	st := *(*C.struct_stat)(unsafe.Pointer(statstruct)) //Casting unsafe pointer to C.struct_stat
	logger.Infoln("Casting unsafe.Pointer(stat) to *(*C.struct_stat)")
	defer C.free(unsafe.Pointer(path)) //free the C.CString that is created in heap.
	C.stat(path, &st)
	uid := st.st_uid
	gid := st.st_gid
	if int(uid) != current_uid || int(gid) != current_gid {
		fmt.Printf("%s not owned by uid=%d(%s) or gid=%d(%s). Do chown %d:%d %s\n", filename, current_uid, current_username, current_gid, current_groupname, current_uid, current_gid, filename)
		logger.Infof("%s not owned by uid=%d(%s) or gid=%d(%s). Do chown %d:%d %s\n", filename, current_uid, current_username, current_gid, current_groupname, current_uid, current_gid, filename)
		os.Exit(1)
	}
	if st.st_mode & C.S_IRGRP > 0 || st.st_mode & C.S_IWGRP > 0 || st.st_mode & C.S_IXGRP > 0 ||
		st.st_mode & C.S_IROTH > 0 || st.st_mode & C.S_IWOTH > 0 || st.st_mode & C.S_IXOTH > 0 {
		fmt.Println(filename, "file permission too broad, make it non-readable to groups and others.")
		logger.Infoln(filename, "file permission too broad, make it non-readable to groups and others.")
		os.Exit(1)
	}
	logger.Infoln("File permission looks secure")
}