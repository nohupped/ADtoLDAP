package gosyncmodules

func CheckForError(e error) {
	if e != nil {
		logger.Errorln(e)
		panic(e)
	}
}
