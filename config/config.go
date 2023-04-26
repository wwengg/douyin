package config

import "os"

// Get returns environment variable returning default value if empty
func Get(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		value = defaultValue
	}

	return value
}
