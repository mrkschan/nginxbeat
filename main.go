package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
)

var Version = "1.0.0-beta1"
var Name = "nginxbeat"

func main() {
	nb := &Nginxbeat{}

	b := beat.NewBeat(Name, Version, nb)

	b.CommandLineSetup()

	b.LoadConfig()
	err := nb.Config(b)
	if err != nil {
		logp.Critical("Config error: %v", err)
		os.Exit(1)
	}

	b.Run()
}
