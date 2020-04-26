package syncer

// CheckForError is just a helper function to check for error and of not nil, logs the error and panic.
func CheckForError(e error) {
	if e != nil {
		logger.Errorln(e)
		panic(e)
	}
}
