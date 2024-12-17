package utils

import "strconv"

func IntSliceToString(data []int, separator string) string {
	var result string

	for i, v := range data {
		f := strconv.Itoa(v)

		if i == 0 {
			result += f
		} else {
			result += separator + f
		}
	}

	return result
}
