package utils

func Any[K comparable](arr []K, predicate func(K) bool) bool {
	for _, v := range arr {
		if predicate(v) {
			return true
		}
	}
	return false
}
