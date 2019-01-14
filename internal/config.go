package internal

import (
	"log"
	"os"
)

// GetEnv retrieve environment property and print it
func GetEnv(item string) string {
	value := os.Getenv(item)
	log.Printf("using env: %s=%s", item, value)
	return value
}
