package beat

type NginxConfig struct {
	// URLs to Nginx status page.
	// Defaults to []string{"http://127.0.0.1/status#stub"}.
	URLs []string

	// Period defines how often to read status in seconds.
	// Defaults to 1 second.
	Period *int64
}

type ConfigSettings struct {
	Input NginxConfig
}
