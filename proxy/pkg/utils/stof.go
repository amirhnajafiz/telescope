package utils

import "fmt"

// StrToFloat64 converts a string to a float64. It returns an error if the conversion fails.
func StrToFloat64(str string) (float64, error) {
	var result float64

	_, err := fmt.Sscanf(str, "%f", &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
