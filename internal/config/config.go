package config

import (
	"os"
)

var (
	Url string
)

func Getenv(key string, defaultKey string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultKey
	}

	return value
}
