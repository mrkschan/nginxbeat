package main

import (
	"github.com/elastic/libbeat/beat"
)

// Nginxbeat implements Beater interface and sends Nginx status using libbeat.
type Nginxbeat struct {
	config ConfigSettings
}

// Config Nginxbeat according to nginxbeat.yml.
func (nb *Nginxbeat) Config(b *beat.Beat) error {
}

// Setup Nginxbeat.
func (nb *Nginxbeat) Setup(b *beat.Beat) error {
}

// Run Nginxbeat.
func (nb *Nginxbeat) Run(b *beat.Beat) error {
}

// Cleanup Nginxbeat.
func (nb *Nginxbeat) Cleanup(b *beat.Beat) error {
}

// Stop Nginxbeat.
func (nb *Nginxbeat) Stop(b *beat.Beat) error {
}
