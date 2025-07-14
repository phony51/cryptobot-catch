package utils

func RemoveDuplicate[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	var list []T
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}
