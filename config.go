package main

type NginxConfig struct {
	// URL to Nginx status page.
	// Defaults to "http://127.0.0.1:8080/status".
	URL string
	// Format of the status page, either "stub" or "json".
	// Use "stub" for Nginx stub status page.
	// Use "json" for Nginx Plus status page.
	// Defaults to "stub".
	Format string

	// Period defines how often to read status in seconds.
	// Defaults to 1 second.
	Period *int64
}

type ConfigSettings struct {
	Input NginxConfig
}
