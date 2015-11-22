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

	p.client.PublishEvent(common.MapStr{
		"@timestamp": now,
		"type":       "nginx",
		"nginx":      s,
	})

	for _, i := range zones {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": now,
			"type":       "zone",
			"zone":       i,
		})
	}

	for _, i := range upstreams {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": now,
			"type":       "upstream",
			"upstream":   i,
		})
	}

	for _, i := range caches {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": now,
			"type":       "cache",
			"cache":      i,
		})
	}

	for _, i := range tcpzones {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": now,
			"type":       "tcpzone",
			"tcpzone":    i,
		})
	}

	for _, i := range tcpupstreams {
		p.client.PublishEvent(common.MapStr{
			"@timestamp":  now,
			"type":        "tcpupstream",
			"tcpupstream": i,
		})
	}
}
