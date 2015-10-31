package parser

// Parser parses status from Nginx status module.
type Parser interface {
	// Parse status from the given url.
	Parse(url string) (map[string]interface{}, error)
}
