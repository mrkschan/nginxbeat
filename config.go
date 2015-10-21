package main

type NginxConfig struct {
	// URL to Nginx status page.
	// Defaults to "http://127.0.0.1:8080/status".
	URL string

	// Period defines how often to read status in seconds.
	// Defaults to 1 second.
	Period *int64
}

type ConfigSettings struct {
	Input NginxConfig
}
