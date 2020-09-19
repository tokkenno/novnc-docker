package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ServerConfig struct {
	Name  string `json:"name"`
	Host  string `json:"host"`
	Proxy string `json:"proxy"`
	Port  int    `json:"port"`
}

type Config struct {
	Port       int            `json:"port"`
	ClientPath string         `json:"client_path"`
	NoVNCPath  string         `json:"novnc_path"`
	Servers    []ServerConfig `json:"servers"`
}

func ReadConfigFile() (Config, error) {
	var config Config
	configPath := "manager.json"

	if _, err := os.Stat(configPath); err != nil {
		configPath, _ = os.UserConfigDir()
		configPath = path.Join(configPath, "novnc", "manager.json")

		if _, err := os.Stat(configPath); err != nil {
			configPath = path.Join("/etc", "novnc", "manager.json")

			if _, err := os.Stat(configPath); err != nil {
				return config, err
			}
		}
	}

	if configFile, err := os.Open(configPath); err == nil {
		defer configFile.Close()

		log.Println(fmt.Sprintf("Configuration file loaded: %s", configPath))

		configByteValue, _ := ioutil.ReadAll(configFile)

		err = json.Unmarshal(configByteValue, &config)

		if err != nil {
			return config, err
		} else {
			return config, nil
		}
	} else {
		log.Fatalf("Error while open config file: %s", err.Error())
		return config, err
	}
}
