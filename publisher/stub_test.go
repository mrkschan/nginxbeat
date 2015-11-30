package publisher

import (
	"encoding/json"
	"testing"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
	"github.com/stretchr/testify/assert"
)

// Mocks publisher.Client.
type StubTestClient struct {
	Channel chan common.MapStr
}

func (c StubTestClient) PublishEvent(event common.MapStr, opts ...publisher.ClientOption) bool {
	c.Channel <- event
	return true
}
func (c StubTestClient) PublishEvents(events []common.MapStr, opts ...publisher.ClientOption) bool {
	for _, event := range events {
		c.Channel <- event
	}
	return true
}

func TestStubPublisher(t *testing.T) {
	c1 := make(chan common.MapStr, 16)
	p1 := NewStubPublisher(&StubTestClient{Channel: c1})

	s1 := map[string]interface{}{}

	p1.Publish(s1, "__SOURCE__")
	assert.Equal(t, 1, len(c1))

	s1e1 := <-c1
	var s1m1 map[string]interface{}
	if err := json.Unmarshal([]byte(s1e1.String()), &s1m1); assert.NoError(t, err) {
		assert.Equal(t, "nginx", s1m1["type"])
		assert.Equal(t, "stub", s1m1["format"])
		assert.Equal(t, "__SOURCE__", s1m1["source"])
	}
}
