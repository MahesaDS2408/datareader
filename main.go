package main

import (
	"datareader/entities"
)

func main() {
	var d *entities.Dependencies
	c := entities.NewYodelConfig()

	// TODO:
	// - Cache config
	// - Glob file sconning
	// - Parsing CSV with Generic thingy
	// - Split out parsed CSV to JSON

	d = &entities.Dependencies{
		Config: c,
	}
	d.Run()
}
