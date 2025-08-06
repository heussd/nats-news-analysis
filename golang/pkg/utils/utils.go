package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(filepath.Join("../..", ".env")); err != nil {
		fmt.Println(err)
	}
}

func Unmarshal(obj any, filename string) (err error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &obj)

	return nil
}

// GetEnv Taken from https://stackoverflow.com//questions/40326540/how-to-assign-default-value-if-env-var-is-empty#answer-45978733
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	fmt.Println(key, "=", value)
	return value
}

func RequireNotEmpty(value, name string) {
	if len(value) == 0 {
		panic(fmt.Sprintf("%s is empty", name))
	}
}
