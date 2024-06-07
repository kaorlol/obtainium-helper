package utils

func Filter[TYPE any](data []TYPE, f func(TYPE) bool) []TYPE {
	var result []TYPE
	for _, d := range data {
		if f(d) {
			result = append(result, d)
		}
	}
	return result
}
