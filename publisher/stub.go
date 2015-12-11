package publisher

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"
)

// StubPublisher is a Publisher that publishes Nginx Stub status.
type StubPublisher struct {
	client publisher.Client
}

// NewStubPublisher constructs a new StubPublisher.
func NewStubPublisher(c publisher.Client) *StubPublisher {
	return &StubPublisher{client: c}
}

// Publish Nginx Stub status.
func (p *StubPublisher) Publish(s map[string]interface{}, source string) {
	p.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "stub",
		"source":     source,
		"stub":       s,
	})
}
