package application

import (
	"encoding/json"
	"os"
)

type Config struct{
	Port int `json:"port"`
	Language string `json:"language"`
	Log_file string `json:"log_file"`
}

func LoadConfig(filePath string) (*Config, error) {
	logger := CreateLogger("./log/server.log")
	logger.Println("Loading configuration from", filePath)
	configFile, err := os.Open(filePath)
	if err != nil {
		logger.Println("Error opening config file:", err)
		return nil, err
	}
	defer configFile.Close()

	config := &Config{}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		logger.Println("Error decoding JSON:", err)
		return nil, err
	}

	return config, nil
}
