package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

type Client struct {
	Address string
	Port    int
}

type DatabaseConfig struct {
	Path string
}

type Config struct {
	MPD               Client `toml:"Client"`
	MusicDirectory    string
	PlaylistDirectory string
	LyricsDirectory   string
	DbConfig          DatabaseConfig `toml:"Database"`
}

func returnConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		if homeDir := os.Getenv("HOME"); homeDir == "" {
			configDir = "~/.config/"
		} else {
			configDir = homeDir + "/" + ".config/"
		}
	}
	return configDir + "innocent/"
}

func returnConfigFile() string {
	return returnConfigDir() + "config.toml"
}

func ReadConfig() *Config {
	var config Config
	dumpDefaultConfig()
	configFile := returnConfigFile()
	fileContent, _ := ioutil.ReadFile(configFile)
	toml.Unmarshal(fileContent, &config)
	return &config
}

func checkIfConfigExists() bool {
	configFile := returnConfigFile()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func getDefaultConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Could not fetch the current user's home directory")
	}
	return Config{
		MusicDirectory:    homeDir + "/Music",
		PlaylistDirectory: homeDir + "/.mpd/playlists",
		LyricsDirectory:   homeDir + "/.mpd/lyrics",
		MPD: Client{
			Address: "localhost",
			Port:    6600,
		},
		DbConfig: DatabaseConfig{
			Path: homeDir + "/.innocent/db.sqlite",
		},
	}
}

func dumpDefaultConfig() {
	if !checkIfConfigExists() {
		config := getDefaultConfig()
		os.MkdirAll(returnConfigDir(), 0766)
		bytesConfig, err := toml.Marshal(config)
		if err != nil {
			log.Fatalln(err)
		}
		ioutil.WriteFile(returnConfigFile(), bytesConfig, 0744)
	}
}
