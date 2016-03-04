// +build unit

package collector

import (
	"os"
	"testing"

	"net"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	// It should hit SOCKS5 proxy.
	p1, _ := net.Listen("tcp", "127.0.0.1:0")
	defer p1.Close()

	pr1 := false
	go func() {
		c, err := p1.Accept()
		if err != nil {
			return
		}
		defer c.Close()

		pr1 = true
	}()

	hr1 := false
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hr1 = true
	}))
	defer ts1.Close()

	os.Setenv("all_proxy", "socks5://"+p1.Addr().String())
	c1 := HTTPClient()
	c1.Get(ts1.URL)
	os.Setenv("all_proxy", "")

	assert.Equal(t, pr1, true)
	assert.Equal(t, hr1, false)

	// It should hit HTTP server directly.
	p2, _ := net.Listen("tcp", "127.0.0.1:0")
	defer p2.Close()

	pr2 := false
	go func() {
		c, err := p2.Accept()
		if err != nil {
			return
		}
		defer c.Close()

		pr2 = true
	}()

	hr2 := false
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hr2 = true
	}))
	defer ts2.Close()

	c2 := HTTPClient()
	c2.Get(ts2.URL)

	assert.Equal(t, pr2, false)
	assert.Equal(t, hr2, true)
}
