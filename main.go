package main

import (
	"github.com/codegangsta/martini"//,
	//"github.com/gorilla/sessions",
	//"github.com/hoisie/mustache"
	// Not used yet, but they will be
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("public"))

	// It will change in the future
	m.Get("/", func() string {
		return "Hello world"
	})

	m.Run()
}