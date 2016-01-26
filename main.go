package main

import (
	"github.com/elastic/beats/libbeat/beat"

	nginxbeat "github.com/mrkschan/nginxbeat/beat"
)

var Version = "1.0.0-beta1"
var Name = "nginxbeat"

func main() {
	beat.Run(Name, Version, nginxbeat.New())
}
