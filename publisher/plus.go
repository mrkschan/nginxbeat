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
	p.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "nginx",
		"nginx":      s,
	})
}
