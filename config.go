package main

import (
	"encoding/json"
	"os"
)

// Config in JSON format
type Config struct {
	ConnectionString string
	SecretKey        string
	Port             string
	SessionExpires   int
}

func Load(path string) *Config {
	reader, _ := os.Open("config.json")
	decoder := json.NewDecoder(reader)
	config := &Config{}
	decoder.Decode(&config)

	return config
}
