package main

func main() {
	config := Load("config.json")

	DB := &Connection{}

	// Init database
	DB.initDatabase(config.ConnectionString)

	Start(DB, config)
}
