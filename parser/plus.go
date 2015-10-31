package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PlusParser is a Parser that parse Nginx Plus status page.
type PlusParser struct {
}

// NewPlusParser constructs a new PlusParser.
func NewPlusParser() Parser {
	return &PlusParser{}
}

// Parse Nginx Plus status from given url.
func (p *PlusParser) Parse(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	// Nginx Plus status is expected to be a JSON document.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	return data, nil
}
