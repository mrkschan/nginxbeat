// +build integration

package collector

import (
	"testing"

	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestStubCollectorIntegration(t *testing.T) {
	// It should report stats of Nginx running inside Docker container.
	c1 := &StubCollector{}
	u1, _ := url.Parse("http://nginx-19:8080/status")
	s1, _ := c1.Collect(*u1)

	assert.Equal(t, s1["active"], 1)
	assert.Equal(t, s1["accepts"], 1)
	assert.Equal(t, s1["handled"], 1)
	assert.Equal(t, s1["dropped"], 0)
	assert.Equal(t, s1["requests"], 1)
	assert.Equal(t, s1["current"], 1)

	// It should report accumlated stats.
	u11, _ := url.Parse("http://nginx-19:8080/status")
	s11, _ := c1.Collect(*u11)

	assert.Equal(t, s11["active"], 1)
	assert.Equal(t, s11["accepts"], 1)
	assert.Equal(t, s11["handled"], 1)
	assert.Equal(t, s11["dropped"], 0)
	assert.Equal(t, s11["requests"], 2)
	assert.Equal(t, s11["current"], 1)
}
