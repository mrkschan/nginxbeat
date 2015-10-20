package main

import (
	"github.com/elastic/libbeat/beat"
)

var Version = "1.0.0-beta1"
var Name = "nginxbeat"

func main() {
	nb := &Nginxbeat{}

	b := beat.NewBeat(Name, Version, nb)

	b.CommandLineSetup()

	b.LoadConfig()
	nb.Config(b)

	b.Run()
}
