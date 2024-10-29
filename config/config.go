package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigSchema struct {
	Address string `json:"address"`
}

func readConfigFile() []byte {
	data, readConfigFileE := os.ReadFile("config.json")
	if readConfigFileE != nil && readConfigFileE == os.ErrNotExist {
		newConfigFile, createConfigFileE := os.Create("config.json")
		if createConfigFileE != nil {
			log.Fatal(createConfigFileE)
		}
		newConfigFile.Write([]byte(`{"address": ":25201"}`))
		newConfigFile.Close()

		data, _ = os.ReadFile("config.json")
	}

	return data
}

func GetConfig() *ConfigSchema {
	configData := readConfigFile()
	config := ConfigSchema{}
	json.Unmarshal(configData, &config)
	return &config
}
