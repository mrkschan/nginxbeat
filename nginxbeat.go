package main

import (
	"net/url"
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/logp"
)

// Nginxbeat implements Beater interface and sends Nginx status using libbeat.
type Nginxbeat struct {
	// NbConfig holds configurations of Nginxbeat parsed by libbeat.
	NbConfig ConfigSettings

	url    *url.URL
	period time.Duration
}

// Config Nginxbeat according to nginxbeat.yml.
func (nb *Nginxbeat) Config(b *beat.Beat) error {
	err := cfgfile.Read(&nb.NbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	var u string
	if nb.NbConfig.Input.URL != "" {
		u = nb.NbConfig.Input.URL
	} else {
		u = "http://127.0.0.1/status"
	}
	nb.url, err = url.Parse(u)
	if err != nil {
		logp.Err("Invalid Nginx status page: %v", err)
		return err
	}

	if nb.NbConfig.Input.Period != nil {
		nb.period = time.Duration(*nb.NbConfig.Input.Period) * time.Second
	} else {
		nb.period = 1 * time.Second
	}

	logp.Debug("nginxbeat", "Init nginxbeat")
	logp.Debug("nginxbeat", "Watch %v\n", nb.url)
	logp.Debug("nginxbeat", "Period %v\n", nb.period)

	return nil
}

// Setup Nginxbeat.
func (nb *Nginxbeat) Setup(b *beat.Beat) error {
	return nil
}

// Run Nginxbeat.
func (nb *Nginxbeat) Run(b *beat.Beat) error {
	return nil
}

// Cleanup Nginxbeat.
func (nb *Nginxbeat) Cleanup(b *beat.Beat) error {
	return nil
}

// Stop Nginxbeat.
func (nb *Nginxbeat) Stop() {
}
