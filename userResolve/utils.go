package userResolve

func findString(array []string, string string) int {
	result := len(array)
	for index, item := range array {
		if item == string {
			result = index
			break
		}
	}
	return result
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
