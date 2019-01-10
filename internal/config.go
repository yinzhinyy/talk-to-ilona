package internal

import (
	"log"
	"os"
)

func GetEnv(item string) string {
	value := os.Getenv(item)
	log.Printf("using env: %s=%s", item, value)
	return value
}
