package genericutils

import (
	"fmt"
	"strconv"
)

func NormalizeFloatValues(values [][]any, length int) []float64 {
	normalizedFloatValues := make([]float64, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			value, err := strconv.ParseFloat(fmt.Sprintf("%v", values[i][0]), 64)
			if err == nil {
				normalizedFloatValues[i] = value
			}
		}
	}

	return normalizedFloatValues
}
