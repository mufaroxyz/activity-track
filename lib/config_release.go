//go:build !debug

package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const DEBUG = 0

var config map[string]string

func InitConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting user home directory")
	}

	configFilePath := filepath.Join(homeDir, ".config", "activity-track.json")
	file, err := os.Open(configFilePath)
	if err != nil {
		panic("Error opening config file: missing file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic("Error decoding config file: invalid JSON")
	}
}

func getEnv(key string) string {
	value, exists := config[key]
	if !exists {
		panic("Key not found in config file: " + key)
	}
	return value
}
