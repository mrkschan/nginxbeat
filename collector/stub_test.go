// +build unit

package collector

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestStubCollector(t *testing.T) {
	// It should report stats.
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 19")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts1.Close()

	c1 := NewStubCollector()
	u1, _ := url.Parse(ts1.URL)
	s1, _ := c1.Collect(*u1)

	assert.Equal(t, s1["active"], 1)
	assert.Equal(t, s1["accepts"], 8)
	assert.Equal(t, s1["handled"], 7)
	assert.Equal(t, s1["dropped"], 1)
	assert.Equal(t, s1["requests"], 19)
	assert.Equal(t, s1["current"], 19)
	assert.Equal(t, s1["reading"], 0)
	assert.Equal(t, s1["writing"], 1)
	assert.Equal(t, s1["waiting"], 2)

	// It should report accumlated stats.
	ts11 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 23")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts11.Close()

	u11, _ := url.Parse(ts11.URL)
	s11, _ := c1.Collect(*u11)

	assert.Equal(t, s11["active"], 1)
	assert.Equal(t, s11["accepts"], 8)
	assert.Equal(t, s11["handled"], 7)
	assert.Equal(t, s11["dropped"], 1)
	assert.Equal(t, s11["requests"], 23)
	assert.Equal(t, s11["current"], 4)
	assert.Equal(t, s11["reading"], 0)
	assert.Equal(t, s11["writing"], 1)
	assert.Equal(t, s11["waiting"], 2)

	// It should handle missing active connections.
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 19")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts2.Close()

	c2 := NewStubCollector()
	u2, _ := url.Parse(ts2.URL)
	s2, _ := c2.Collect(*u2)

	assert.Equal(t, s2["active"], -1)
	assert.Equal(t, s2["accepts"], 8)
	assert.Equal(t, s2["handled"], 7)
	assert.Equal(t, s2["dropped"], 1)
	assert.Equal(t, s2["requests"], 19)
	assert.Equal(t, s2["current"], 19)
	assert.Equal(t, s2["reading"], 0)
	assert.Equal(t, s2["writing"], 1)
	assert.Equal(t, s2["waiting"], 2)

	// It should handle missing request stats.
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts3.Close()

	c3 := NewStubCollector()
	u3, _ := url.Parse(ts3.URL)
	s3, _ := c3.Collect(*u3)

	assert.Equal(t, s3["active"], 1)
	assert.Equal(t, s3["accepts"], -1)
	assert.Equal(t, s3["handled"], -1)
	assert.Equal(t, s3["dropped"], -1)
	assert.Equal(t, s3["requests"], -1)
	assert.Equal(t, s3["current"], -1)
	assert.Equal(t, s3["reading"], 0)
	assert.Equal(t, s3["writing"], 1)
	assert.Equal(t, s3["waiting"], 2)

	// It should handle missing connection stats.
	ts4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 19")
		fmt.Fprintln(w, "")
	}))
	defer ts4.Close()

	c4 := NewStubCollector()
	u4, _ := url.Parse(ts4.URL)
	s4, _ := c4.Collect(*u4)

	assert.Equal(t, s4["active"], 1)
	assert.Equal(t, s4["accepts"], 8)
	assert.Equal(t, s4["handled"], 7)
	assert.Equal(t, s4["dropped"], 1)
	assert.Equal(t, s4["requests"], 19)
	assert.Equal(t, s4["current"], 19)
	assert.Equal(t, s4["reading"], -1)
	assert.Equal(t, s4["writing"], -1)
	assert.Equal(t, s4["waiting"], -1)

	// It should report unexpected status code.
	ts5 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failed", http.StatusInternalServerError)
	}))
	defer ts5.Close()

	c5 := NewStubCollector()
	u5, _ := url.Parse(ts5.URL)
	s5, e5 := c5.Collect(*u5)

	assert.Nil(t, s5)
	assert.EqualError(t, e5, "HTTP500 Internal Server Error")
}
