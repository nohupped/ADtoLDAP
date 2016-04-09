package gosyncmodules

func StringInSlice(checkfor string, checkin []string) bool {
	for _, i := range checkin {
		if i == checkfor {
			return true
		}
	}
	return false
}
