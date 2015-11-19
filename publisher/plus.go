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

	stream := s["stream"]
	delete(s, "stream")

	p.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "nginx",
		"nginx":      s,
	})

	for _, i := range zones {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       "zone",
			"zone":       i,
		})
	}

	for _, i := range upstreams {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       "upstream",
			"upstream":   i,
		})
	}

	for _, i := range caches {
		p.client.PublishEvent(common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       "cache",
			"cache":      i,
		})
	}

	p.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "stream",
		"stream":     stream,
	})
}
