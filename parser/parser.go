package parser

import (
	"net/url"
)

// Parser parses status from Nginx status module.
type Parser interface {
	// Parse status from the given url.
	Parse(u url.URL) (map[string]interface{}, error)
}
