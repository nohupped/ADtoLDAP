package gosyncmodules

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// CheckPerm checks permission of the passed filename(string) and panics if group or others have > read permissions.
func CheckPerm(filename string) {
	currentUID := os.Getuid()
	currentGID := os.Getgid()
	username, err := user.LookupId(strconv.Itoa(currentUID))
	var currentUserName string
	var currentGroupName string
	if err != nil {
		currentUserName = "<Couldn't lookup username>"
	}
	groupname, err := user.LookupGroupId(strconv.Itoa(currentGID))
	if err != nil {
		currentGroupName = "<Couldn't lookup groupname>"

	}
	currentUserName = username.Username
	currentGroupName = groupname.Name

	fileStat, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	filestatSys := fileStat.Sys().(*syscall.Stat_t)
	uid := filestatSys.Uid
	gid := filestatSys.Gid
	if int(uid) != currentUID || int(gid) != currentGID {
		fmt.Printf("%s not owned by uid=%d(%s) or gid=%d(%s). Do chown %d:%d %s\n", filename, currentUID, currentUserName, currentGID, currentGroupName, currentUID, currentGID, filename)
		logger.Errorf("%s not owned by uid=%d(%s) or gid=%d(%s). Do chown %d:%d %s\n", filename, currentUID, currentUserName, currentGID, currentGroupName, currentUID, currentGID, filename)
		os.Exit(1)
	}
	if filestatSys.Mode&syscall.S_IRGRP > 0 || filestatSys.Mode&syscall.S_IWGRP > 0 || filestatSys.Mode&syscall.S_IXGRP > 0 ||
		filestatSys.Mode&syscall.S_IROTH > 0 || filestatSys.Mode&syscall.S_IWOTH > 0 || filestatSys.Mode&syscall.S_IXOTH > 0 {
		fmt.Println(filename, "file permission too broad, make it non-readable to groups and others.")
		logger.Errorln(filename, "file permission too broad, make it non-readable to groups and others.")
		os.Exit(1)
	}
	logger.Infoln("File permission looks secure")
}
