package publisher

import (
	"time"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
)

// PlusPublisher is a Publisher that publishes Nginx Plus status.
type PlusPublisher struct {
	client publisher.Client
}

// NewPlusPublisher constructs a new PlusPublisher.
func NewPlusPublisher(c publisher.Client) *PlusPublisher {
	return &PlusPublisher{client: c}
}

// Publish Nginx Plus status.
func (p *PlusPublisher) Publish(s map[string]interface{}) {
	const format = "plus"

	version := s["version"]
	nginxVersion := s["nginx_version"]

	zones := s["server_zones"].([]interface{})
	delete(s, "server_zones")

	upstreams := s["upstreams"].([]interface{})
	delete(s, "upstreams")

	caches := s["caches"].([]interface{})
	delete(s, "caches")

	stream := s["stream"].(map[string]interface{})
	delete(s, "stream")

	tcpzones := stream["server_zones"].([]interface{})
	tcpupstreams := stream["upstreams"].([]interface{})

	now := common.Time(time.Now())

	buf := []common.MapStr{}

	buf = append(buf, common.MapStr{
		"@timestamp": now,
		"type":       "nginx",
		"format":     format,
		"nginx":      s,
	})

	for _, i := range zones {
		m := i.(map[string]interface{})
		m["version"] = version
		m["nginx_version"] = nginxVersion

		buf = append(buf, common.MapStr{
			"@timestamp": now,
			"type":       "zone",
			"format":     format,
			"zone":       m,
		})
	}

	for _, i := range upstreams {
		m := i.(map[string]interface{})
		m["version"] = version
		m["nginx_version"] = nginxVersion

		buf = append(buf, common.MapStr{
			"@timestamp": now,
			"type":       "upstream",
			"format":     format,
			"upstream":   m,
		})
	}

	for _, i := range caches {
		m := i.(map[string]interface{})
		m["version"] = version
		m["nginx_version"] = nginxVersion

		buf = append(buf, common.MapStr{
			"@timestamp": now,
			"type":       "cache",
			"format":     format,
			"cache":      m,
		})
	}

	for _, i := range tcpzones {
		m := i.(map[string]interface{})
		m["version"] = version
		m["nginx_version"] = nginxVersion

		buf = append(buf, common.MapStr{
			"@timestamp": now,
			"type":       "tcpzone",
			"format":     format,
			"tcpzone":    m,
		})
	}

	for _, i := range tcpupstreams {
		m := i.(map[string]interface{})
		m["version"] = version
		m["nginx_version"] = nginxVersion

		buf = append(buf, common.MapStr{
			"@timestamp":  now,
			"type":        "tcpupstream",
			"format":      format,
			"tcpupstream": m,
		})
	}

	p.client.PublishEvents(buf)
}
