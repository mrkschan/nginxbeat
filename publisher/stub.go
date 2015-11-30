package publisher

import (
	"time"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
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
func (p *StubPublisher) Publish(s map[string]interface{}) {
	const format = "stub"

	s["format"] = format
	p.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "nginx",
		"nginx":      s,
	})
}
