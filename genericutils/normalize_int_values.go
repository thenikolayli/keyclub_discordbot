package genericutils

import (
	"fmt"
	"strconv"
)

func NormalizeIntValues(values [][]any, length int) []int {
	normalizedIntValues := make([]int, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			value, err := strconv.ParseInt(fmt.Sprintf("%v", values[i][0]), 0, 64)
			if err == nil {
				normalizedIntValues[i] = int(value)
			}
		}
	}

	return normalizedIntValues
}
