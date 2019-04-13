package gosyncmodules

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// Constants S_IRGRP, S_IWGRP, S_IXGRP, S_IROTH, S_IWOTH and S_IXOTH values are taken straight from /stat.h.
// Ref: https://github.com/torvalds/linux/blob/master/include/uapi/linux/stat.h
const (
	IRGRP = 0000040
	IWGRP = 0000020
	IXGRP = 0000010
	IROTH = 0000004
	IWOTH = 0000002
	IXOTH = 0000001
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
	if filestatSys.Mode&IRGRP > 0 || filestatSys.Mode&IWGRP > 0 || filestatSys.Mode&IXGRP > 0 ||
	   filestatSys.Mode&IROTH > 0 || filestatSys.Mode&IWOTH > 0 || filestatSys.Mode&IXOTH > 0 {
		fmt.Println(filename, "file permission too broad, make it non-readable to groups and others.")
		logger.Errorln(filename, "file permission too broad, make it non-readable to groups and others.")
		os.Exit(1)
	}
	logger.Infoln("File permission looks secure")
}

