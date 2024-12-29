package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
	filepath        string
}

func Read() (*Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filepath := homedir + "/" + configFileName

	f, err := os.OpenFile(filepath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}
	c.filepath = filepath

	return &c, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	f, err := os.OpenFile(c.filepath, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "\t")
	return e.Encode(c)
}
