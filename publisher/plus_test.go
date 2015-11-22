package publisher

import (
	"encoding/json"
	"testing"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
	"github.com/stretchr/testify/assert"
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
		"version": 6,
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
		"stream": map[string]interface{}{},
	}

	p1.Publish(s1)
	assert.Equal(t, 5, len(c1))

	s1e1 := <-c1
	var s1m1 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e1.String()), &s1m1); assert.NoError(t, err) {
		assert.Equal(t, "nginx", s1m1["type"])
	}

	s1e2 := <-c1
	var s1m2 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e2.String()), &s1m2); assert.NoError(t, err) {
		assert.Equal(t, "zone", s1m2["type"])
	}

	s1e3 := <-c1
	var s1m3 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e3.String()), &s1m3); assert.NoError(t, err) {
		assert.Equal(t, "upstream", s1m3["type"])
	}

	s1e4 := <-c1
	var s1m4 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e4.String()), &s1m4); assert.NoError(t, err) {
		assert.Equal(t, "cache", s1m4["type"])
	}

	s1e5 := <-c1
	var s1m5 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e5.String()), &s1m5); assert.NoError(t, err) {
		assert.Equal(t, "stream", s1m5["type"])
	}
}
