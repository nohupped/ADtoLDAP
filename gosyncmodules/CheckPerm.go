package gosyncmodules

import (
	"fmt"
	"unsafe"
	"os"
)
//#include <sys/stat.h>
//#include <stdlib.h>
import "C"

func CheckPerm(filename string) {
	Info.Println("using cgo to perform security check on ", filename)
	statt := C.stat //stat struct from C
	Info.Println("Initiated stat struct")
	path := C.CString(filename)
	Info.Println("Converted native string to C.CString")
	st := *(*C.struct_stat)(unsafe.Pointer(statt)) //Casting unsafe pointer to C.struct_stat
	Info.Println("Casting unsafe.Pointer(stat) to *(*C.struct_stat)")
	defer C.free(unsafe.Pointer(path)) //free
	C.stat(path, &st)
	if st.st_mode & C.S_IRGRP > 0 || st.st_mode & C.S_IWGRP > 0 || st.st_mode & C.S_IXGRP > 0 || st.st_mode & C.S_IROTH > 0 || st.st_mode & C.S_IWOTH > 0 || st.st_mode & C.S_IXOTH > 0 {
		fmt.Println("file permission too broad, make it non-readable to groups and others.")
		os.Exit(1)
	}
	Info.Println("File permission looks secure")
}