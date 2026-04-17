package genericutils

import "fmt"

// the following functions create normalized lists with blanks from the originals
func NormalizeStringValues(values [][]any, length int) []string {
	normalizedStringValues := make([]string, length)

	for i := range values {
		if len(values[i]) != 0 { // if the cell isn't blank
			normalizedStringValues[i] = fmt.Sprintf("%v", values[i][0])
		}
	}

	return normalizedStringValues
}
