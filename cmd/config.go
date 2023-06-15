package main

import (
	"ege/app/ege"
	"ege/app/vk"
	"encoding/json"
	"os"
)

const fileName = "config.json"

// Config reads the configuration file.
type Config struct {
	IsFilled bool       `json:"is_filled"`
	Ege      ege.Config `json:"ege"`
	VK       vk.Config  `json:"vk"`
}

// NewConfig initializes the configuration file.
func newConfig() Config {
	return Config{
		Ege: ege.NewConfig(),
		VK:  vk.NewConfig(),
	}
}

func OverwriteConfig() error {

	config := newConfig()
	jsonString, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return NewInitializationError(MarshalConfigError, err)
	}

	err = os.WriteFile(fileName, jsonString, os.ModePerm)
	if err != nil {
		return NewInitializationError(WriteConfigError, err)
	}

	return nil
}

// GetConfig from config.json or default config.
func GetConfig() (*Config, error) {

	config := &Config{}

	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewInitializationError(MissingConfigFile, err)
		}
		return nil, NewInitializationError(ReadConfigError, err)
	}

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, NewInitializationError(ReadConfigError, err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, NewInitializationError(UnmarshalConfigError, err)
	}

	if config.IsFilled == false {
		return nil, NewInitializationError(ConfigIsNotFiledError, err)
	}

	return config, nil

}
