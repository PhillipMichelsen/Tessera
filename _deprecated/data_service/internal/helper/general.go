package helper

import (
	"log"
	"strconv"
)

// FormatFloat converts a float64 to a string with a fixed precision
func FormatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

func ParseFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing float value %s: %v", value, err)
		return 0.0
	}
	return result
}
