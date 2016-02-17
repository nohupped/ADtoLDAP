package gosyncmodules

func CheckForError(e error) {
	if e != nil {
		Error.Println(e)
		panic(e)
	}
}
