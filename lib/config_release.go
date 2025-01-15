//go:build !debug

package lib

import (
	"github.com/joho/godotenv"
	"os"
)

func getEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv(key)
}
