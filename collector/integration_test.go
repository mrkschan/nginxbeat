// +build integration

package collector

import (
	"testing"

	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestStubCollectorIntegration(t *testing.T) {
	// It should report stats of Nginx running inside Docker container.
	c1 := NewStubCollector()
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

func TestPlusCollectorIntegration(t *testing.T) {
	// It should report stats of Nginx running on demo.nginx.com.
	c1 := NewPlusCollector()
	u1, _ := url.Parse("http://demo.nginx.com/status")
	s1, _ := c1.Collect(*u1)

	sver1 := s1["version"]
	assert.True(t, sver1.(int) >= 6, "Status version should be >=6.")
	sconns1 := s1["connections"].(map[string]interface{})
	assert.True(t, sconns1["active"].(int) > 0, "Active connections should be greater than 0.")
	sreqs1 := s1["requests"].(map[string]interface{})
	assert.True(t, sreqs1["current"].(int) > 0, "No. of requests should be greater than 0.")

	szones1 := s1["server_zones"].([]interface{})
	assert.True(t, len(szones1) > 0, "No. of server zones should be greater than 0.")
	supstreams1 := s1["upstreams"].([]interface{})
	assert.True(t, len(supstreams1) > 0, "No. of upstreams should be greater than 0.")
	scaches1 := s1["caches"].([]interface{})
	assert.True(t, len(scaches1) > 0, "No. of cache zones should be greater than 0.")

	sstreams1 := s1["stream"].(map[string]interface{})
	sstreamzones1 := sstreams1["server_zones"].([]interface{})
	sstreamupstreams1 := sstreams1["upstreams"].([]interface{})
	assert.True(t, len(sstreamzones1) > 0, "No. of TCP stream server zones should be greater than 0.")
	assert.True(t, len(sstreamupstreams1) > 0, "No. of TCP stream upstreams should be greater than 0.")
}
