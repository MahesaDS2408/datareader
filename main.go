package main

import (
	"datareader/deliveries"
	"datareader/entities"
)

func main() {
	var d *entities.Dependencies
	c := entities.NewYodelConfig()
	r := deliveries.NewHttpRoute(d)

	// TODO:
	// - Cache config
	// - Glob file sconning
	// - Parsing CSV with Generic thingy
	// - Split out parsed CSV to JSON

	d = &entities.Dependencies{
		Config: c,
		Router: r,
	}
	d.Run()
}
