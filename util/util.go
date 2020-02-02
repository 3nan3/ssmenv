package util

func SliceContains(array []string, key string) bool {
	for _, elem := range array {
		if elem == key {
			return true
		}
	}
	return false
}
