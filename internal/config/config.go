//go:build !debug

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var config map[string]string

func createConfigFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting user home directory")
	}

	configFilePath := filepath.Join(homeDir, ".config", "activity-track.json")
	file, err := os.Create(configFilePath)
	if err != nil {
		panic("Error creating config file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.WriteString(`{
	"CF_API_KEY": "Your cloudflare API key (with D1 edit permission)",
	"CF_ACCOUNT_ID": "Your cloudflare account id",
	"D1_ID": "Your D1 database id"
}`)
	if err != nil {
		panic("Error writing to config file")
	}

	println("Config file created at: " + configFilePath)
	os.Exit(0)
}

func InitConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Error getting user home directory")
	}

	configFilePath := filepath.Join(homeDir, ".config", "activity-track.json")
	file, err := os.Open(configFilePath)
	if err != nil {
		createConfigFile()
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

func GetEnv(key string) string {
	value, exists := config[key]
	if !exists {
		panic("Key not found in config file: " + key)
	}
	return value
}
