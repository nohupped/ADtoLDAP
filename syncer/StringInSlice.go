package syncer

// StringInSlice is a helper function to check if a string exists in a given slice.
func StringInSlice(checkfor string, checkin []string) bool {
	for _, i := range checkin {
		if i == checkfor {
			return true
		}
	}
	return false
}
