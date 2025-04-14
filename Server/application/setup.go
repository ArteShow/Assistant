package application

import (
	"encoding/json"
	"os"
	"log"
)

type Config struct{
	Port int `json:"port"`
	Language string `json:"language"`
	Log_file string `json:"log_file"`
}

func LoadConfig(filePath string) (*Config, error) {
	log.Println("Loading configuration from", filePath)
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening config file:", err)
		return nil, err
	}
	defer configFile.Close()

	config := &Config{}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return nil, err
	}

	return config, nil
}
