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

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, err
	}

	// Convert float to int values.
	status := ftoi(payload)


	return status, nil
}

// ftoi returns a copy of in where float values are casted to int values.
func ftoi(in map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}

	for k, v := range in {
		switch v.(type) {
		case float64:
			vt := v.(float64)
			out[k] = int(vt)
		case map[string]interface{}:
			vt := v.(map[string]interface{})
			out[k] = ftoi(vt)
		case []interface{}:
			vt := v.([]interface{})
			l := len(vt)
			a := make([]interface{}, l)
			for i := 0; i < l; i++ {
				e := vt[i]
				switch e.(type) {
				case float64:
					et := e.(float64)
					a[i] = int(et)
				case map[string]interface{}:
					et := e.(map[string]interface{})
					a[i] = ftoi(et)
				default:
					a[i] = e
				}
			}
			out[k] = a
		default:
			out[k] = v
		}
	}

	return out
}
