package main

import (
	"os"
	"encoding/json"
	"github.com/gorilla/sessions"
)

func main() {
	// Read config
	reader, _ := os.Open("config.json")
	decoder := json.NewDecoder(reader)
	config := &Config{}
	decoder.Decode(&config)

	DB := &Connection{}

	// Init database
	DB.initDatabase(config.ConnectionString)

	// Create a new cookie store
	store := sessions.NewCookieStore([]byte(config.SecretKey))

	Start(DB, store, config)
}
