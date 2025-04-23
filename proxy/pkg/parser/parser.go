package parser

import (
	"fmt"
	"strings"
)

// ParseHeaderToMap parses a header string into a map
func ParseMapToHeader(headerMap map[string]float64) string {
	header := ""

	for key, value := range headerMap {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf("%s:%f", key, value)
	}

	return header
}

// ParseHeaderToMap parses a header string into a map
func ParseHeaderToMap(header string) map[string]float64 {
	headerMap := make(map[string]float64)
	pairs := strings.SplitSeq(header, ",")

	for pari := range pairs {
		kv := strings.Split(pari, ":")
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		if key == "" || value == "" {
			continue
		}

		headerMap[key] = toFloat64(value)
	}

	return headerMap
}

// converts a string to a float64
func toFloat64(str string) float64 {
	var result float64

	_, err := fmt.Sscanf(str, "%f", &result)
	if err != nil {
		return 0
	}

	return result
}
