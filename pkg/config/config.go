package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type Client struct {
	Address string
	Port    int
}

type Config struct {
	MPD               Client `toml:"Client"`
	MusicDirectory    string
	PlaylistDirectory string
	LyricsDirectory   string
}

func returnConfigFile() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		if homeDir := os.Getenv("HOME"); homeDir == "" {
			configDir = "~/.config/"
		} else {
			configDir = homeDir + "/" + ".config/"
		}
	}
	return configDir + "innocent/config.toml"
}

func ReadConfig() *Config {
	var config Config
	configFile := returnConfigFile()
	fileContent, err := ioutil.ReadFile(configFile)
	if err == nil {
		if _, err = toml.Decode(string(fileContent), &config); err != nil {
			return &Config{}
		}
	}
	return &config
}

func checkIfConfigExists() bool {
	configFile := returnConfigFile()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return false
	}
	return true
}
