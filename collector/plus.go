package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// PlusCollector is a Collector that collects Nginx Plus status page.
type PlusCollector struct {
	http *http.Client
}

// NewPlusCollector constructs a new PlusCollector.
func NewPlusCollector() Collector {
	return &PlusCollector{http: HTTPClient()}
}

// Collect Nginx Plus status from given url.
func (c *PlusCollector) Collect(u url.URL) (map[string]interface{}, error) {
	res, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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
	status := Ftoi(payload)
	version := status["version"].(int)

	// Convert 'server_zones' into nested data type instead of
	// object data type with arbitrary keys.
	if version >= 2 {
		mapping := status["server_zones"].(map[string]interface{})
		zones := make([]interface{}, len(mapping))
		i := 0
		for k, v := range mapping {
			vt := v.(map[string]interface{})
			vt["name"] = k
			zones[i] = vt
			i++
		}
		status["server_zones"] = zones
	} else {
		status["server_zones"] = []interface{}{}
	}

	// Convert 'upstreams' into nested data type instead of
	// object data type with arbitrary keys.
	if version >= 4 {
		mapping := status["upstreams"].(map[string]interface{})
		upstreams := make([]interface{}, len(mapping))
		i := 0
		for k, v := range mapping {
			vt := v.(map[string]interface{})
			vt["name"] = k
			upstreams[i] = vt
			i++
		}
		status["upstreams"] = upstreams
	} else {
		status["upstreams"] = []interface{}{}
	}

	// Convert 'caches' into nested data type instead of
	// object data type with arbitrary keys.
	if version >= 2 {
		mapping := status["caches"].(map[string]interface{})
		caches := make([]interface{}, len(mapping))
		i := 0
		for k, v := range mapping {
			vt := v.(map[string]interface{})
			vt["name"] = k
			caches[i] = vt
			i++
		}
		status["caches"] = caches
	} else {
		status["caches"] = []interface{}{}
	}

	// Convert 'stream' into nested data type instead of
	// object data type with arbitrary keys.
	if version >= 5 {
		var (
			mapping map[string]interface{}
			i       int
		)
		stream := status["stream"].(map[string]interface{})

		mapping = stream["server_zones"].(map[string]interface{})
		zones := make([]interface{}, len(mapping))
		i = 0
		for k, v := range mapping {
			vt := v.(map[string]interface{})
			vt["name"] = k
			zones[i] = vt
			i++
		}
		stream["server_zones"] = zones

		mapping = stream["upstreams"].(map[string]interface{})
		upstreams := make([]interface{}, len(mapping))
		i = 0
		for k, v := range mapping {
			vt := v.(map[string]interface{})
			vt["name"] = k
			upstreams[i] = vt
			i++
		}
		stream["upstreams"] = upstreams

		status["stream"] = stream
	} else {
		status["stream"] = map[string]interface{}{
			"server_zones": []interface{}{},
			"upstreams":    []interface{}{},
		}
	}

	return status, nil
}
