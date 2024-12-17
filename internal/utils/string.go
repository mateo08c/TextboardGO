package utils

import (
	"strconv"
	"strings"
)

func StringToIntSlice(data string, separator string) []int {
	var result []int

	for _, v := range strings.Split(data, separator) {
		i, _ := strconv.Atoi(v)
		result = append(result, i)
	}

	return result
}
