package publisher

import (
	"encoding/json"
	"testing"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
	"github.com/stretchr/testify/assert"

	"github.com/mrkschan/nginxbeat/collector"
)

// Mocks publisher.Client.
type PlusTestClient struct {
	Channel chan common.MapStr
}

func (c PlusTestClient) PublishEvent(event common.MapStr, opts ...publisher.ClientOption) bool {
	c.Channel <- event
	return true
}
func (c PlusTestClient) PublishEvents(events []common.MapStr, opts ...publisher.ClientOption) bool {
	for _, event := range events {
		c.Channel <- event
	}
	return true
}

func TestPlusPublisher(t *testing.T) {
	c1 := make(chan common.MapStr, 16)
	p1 := NewPlusPublisher(&PlusTestClient{Channel: c1})

	s1 := map[string]interface{}{
		"version":       6,
		"nginx_version": "1.9.4",
		"server_zones": []interface{}{
			map[string]interface{}{
				"name": "t.nginx.org",
			},
		},
		"upstreams": []interface{}{
			map[string]interface{}{
				"name": "t.nginx.org",
			},
		},
		"caches": []interface{}{
			map[string]interface{}{
				"name": "http_cache",
			},
		},
		"stream": map[string]interface{}{
			"server_zones": []interface{}{
				map[string]interface{}{
					"name": "tcp_zone",
				},
			},
			"upstreams": []interface{}{
				map[string]interface{}{
					"name": "tcp_upstream",
				},
			},
		},
	}

	p1.Publish(s1, "__SOURCE__")
	assert.Equal(t, 6, len(c1))

	s1e1 := <-c1
	var s1m1 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e1.String()), &s1m1); assert.NoError(t, err) {
		_p := collector.Ftoi(s1m1)
		assert.Equal(t, "nginx", _p["type"])
		assert.Equal(t, "plus", _p["format"])
		assert.Equal(t, "__SOURCE__", _p["source"])

		_s := _p["nginx"].(map[string]interface{})
		assert.Equal(t, 6, _s["version"])
		assert.Equal(t, "1.9.4", _s["nginx_version"])
	}

	s1e2 := <-c1
	var s1m2 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e2.String()), &s1m2); assert.NoError(t, err) {
		_p := collector.Ftoi(s1m2)
		assert.Equal(t, "zone", _p["type"])
		assert.Equal(t, "plus", _p["format"])
		assert.Equal(t, "__SOURCE__", _p["source"])

		_s := _p["zone"].(map[string]interface{})
		assert.Equal(t, 6, _s["version"])
		assert.Equal(t, "1.9.4", _s["nginx_version"])
	}

	s1e3 := <-c1
	var s1m3 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e3.String()), &s1m3); assert.NoError(t, err) {
		_p := collector.Ftoi(s1m3)
		assert.Equal(t, "upstream", _p["type"])
		assert.Equal(t, "plus", _p["format"])
		assert.Equal(t, "__SOURCE__", _p["source"])

		_s := _p["upstream"].(map[string]interface{})
		assert.Equal(t, 6, _s["version"])
		assert.Equal(t, "1.9.4", _s["nginx_version"])
	}

	s1e4 := <-c1
	var s1m4 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e4.String()), &s1m4); assert.NoError(t, err) {
		_p := collector.Ftoi(s1m4)
		assert.Equal(t, "cache", _p["type"])
		assert.Equal(t, "plus", _p["format"])
		assert.Equal(t, "__SOURCE__", _p["source"])

		_s := _p["cache"].(map[string]interface{})
		assert.Equal(t, 6, _s["version"])
		assert.Equal(t, "1.9.4", _s["nginx_version"])
	}

	s1e5 := <-c1
	var s1m5 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e5.String()), &s1m5); assert.NoError(t, err) {
		_p := collector.Ftoi(s1m5)
		assert.Equal(t, "tcpzone", _p["type"])
		assert.Equal(t, "plus", _p["format"])
		assert.Equal(t, "__SOURCE__", _p["source"])

		_s := _p["tcpzone"].(map[string]interface{})
		assert.Equal(t, 6, _s["version"])
		assert.Equal(t, "1.9.4", _s["nginx_version"])
	}

	s1e6 := <-c1
	var s1m6 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e6.String()), &s1m6); assert.NoError(t, err) {
		_p := collector.Ftoi(s1m6)
		assert.Equal(t, "tcpupstream", _p["type"])
		assert.Equal(t, "plus", _p["format"])
		assert.Equal(t, "__SOURCE__", _p["source"])

		_s := _p["tcpupstream"].(map[string]interface{})
		assert.Equal(t, 6, _s["version"])
		assert.Equal(t, "1.9.4", _s["nginx_version"])
	}
}
