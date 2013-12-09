package main

// Config in JSON format
type Config struct {
	ConnectionString string
	SecretKey        string
	Port             string
	SessionExpires   int
}
