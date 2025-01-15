//go:build debug

package lib

import (
	"github.com/joho/godotenv"
	"os"
)

const DEBUG = 1

func getEnv(key string) string {
	println("Loading .env file")
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv(key)
}
