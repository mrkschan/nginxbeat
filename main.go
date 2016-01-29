package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	nginxbeat "github.com/mrkschan/nginxbeat/beat"
)

var Version = "1.0.0-beta1"
var Name = "nginxbeat"

func main() {
	err := beat.Run(Name, Version, nginxbeat.New())
	if err != nil {
		os.Exit(1)
	}
}
