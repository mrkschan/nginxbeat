package collector

import (
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
)

// Collector collects status from Nginx status module.
type Collector interface {
	// Collect status from the given url.
	Collect(u url.URL) (map[string]interface{}, error)
}

// HTTPClient returns a net/http client that respects proxy settings from the
// environmnent. Supported environmnent variables:
// "http_proxy", "https_proxy", "all_proxy", and "no_proxy".
func HTTPClient() *http.Client {
	if os.Getenv("all_proxy") != "" {
		return &http.Client{
			Transport: &http.Transport{Dial: proxy.FromEnvironment().Dial},
		}
	}

	return http.DefaultClient
}

// Ftoi returns a copy of input map where float values are casted to int values.
// The conversion is applied to nested maps and arrays as well.
func Ftoi(in map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}

	for k, v := range in {
		switch v.(type) {
		case float64:
			vt := v.(float64)
			out[k] = int(vt)
		case map[string]interface{}:
			vt := v.(map[string]interface{})
			out[k] = Ftoi(vt)
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
					a[i] = Ftoi(et)
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
