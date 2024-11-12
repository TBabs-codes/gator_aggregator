package config

import (
	"encoding/json"
	"os"
)

const pathToConfig = "/Users/thomasbabcock/.gatorconfig.json"

type Config struct {
	DbURL string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {

	file_content, err := os.ReadFile(pathToConfig)

	if err != nil {
		return Config{}, err
	}

	current_config := Config{}
	err = json.Unmarshal(file_content, &current_config)
	if err != nil {
		return Config{}, err
	}

	return current_config, nil
}

func (c Config) SetUser(user string) error {
	c.CurrentUser = user
	return write(c)
}

func write(c Config) error {
	file_bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	
	err = os.WriteFile(pathToConfig, file_bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}