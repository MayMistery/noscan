package utils

// RemoveDuplicateStringArr removes duplicate strings from a slice.
func RemoveDuplicateStringArr(p []string) []string {
	result := make([]string, 0, len(p))
	temp := map[string]struct{}{}
	for _, item := range p {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
