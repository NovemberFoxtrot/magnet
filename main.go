package main

import (
	"os"
	"encoding/json"
	"github.com/gorilla/sessions"
)

func loadConfig(path string) *Config {
	reader, _ := os.Open("config.json")
	decoder := json.NewDecoder(reader)
	config := &Config{}
	decoder.Decode(&config)

	return config
}

func main() {
	config := loadConfig("config.json")

	DB := &Connection{}

	// Init database
	DB.initDatabase(config.ConnectionString)

	Start(DB, config)
}
