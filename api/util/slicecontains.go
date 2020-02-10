package util

func Contains(slice []string, val string) bool {
	for _, i := range slice {
		if i == val {
			return true
		}
	}
	return false
}