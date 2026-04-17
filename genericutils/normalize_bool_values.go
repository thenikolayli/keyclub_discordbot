package genericutils

import "strconv"

func NormalizeBoolValues(values [][]any, length int) []bool {
	normalizedStringValues := make([]bool, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			value, err := strconv.ParseBool(values[i][0].(string))
			if err == nil {
				normalizedStringValues[i] = bool(value)
			}
		}
	}

	return normalizedStringValues
}
