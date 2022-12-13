package utils

import (
	"fmt"
	"os"
)

// GetEnv Taken from https://stackoverflow.com//questions/40326540/how-to-assign-default-value-if-env-var-is-empty#answer-45978733
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	fmt.Println(key, "=", value)
	return value
}
