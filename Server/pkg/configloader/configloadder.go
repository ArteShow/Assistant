package configloader

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Application_port  int    `json:"application_port"`
	Internal_port     int    `json:"internal_port"`
	Language          string `json:"language"`
	Log_file          string `json:"log_file"`
	Database_path     string `json:"database_path"`
	MoneyDataBasePath string `json:"money_database_pathee"`
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
	log.Println("Configuration loaded successfully")
	return config, nil
}

func GetDatabasePath() (string, error) {
	config, err := LoadConfig("Server/configs/config.json")
	return config.Database_path, err
}

func GetMoneyDatabasePath() (string, error) {
	config, err := LoadConfig("Server/configs/config.json")
	return config.MoneyDataBasePath, err
}

func GetApplicationPort() (int, error) {
	config, err := LoadConfig("Server/configs/config.json")
	return config.Application_port, err
}

func GetInternalPort() (int, error) {
	config, err := LoadConfig("Server/configs/config.json")
	return config.Internal_port, err
}

func GetLanguage() (string, error) {
	config, err := LoadConfig("Server/configs/config.json")
	return config.Language, err
}
