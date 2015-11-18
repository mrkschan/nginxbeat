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
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 19")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts1.Close()

	c1 := &StubCollector{}
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

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 19")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts2.Close()

	c2 := &StubCollector{}
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

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts3.Close()

	c3 := &StubCollector{}
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

	ts4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 19")
		fmt.Fprintln(w, "")
	}))
	defer ts4.Close()

	c4 := &StubCollector{}
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

	ts41 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Active connections: 1")
		fmt.Fprintln(w, "server accepts handled requests")
		fmt.Fprintln(w, " 8 7 23")
		fmt.Fprintln(w, "Reading: 0 Writing: 1 Waiting: 2")
	}))
	defer ts41.Close()

	u41, _ := url.Parse(ts41.URL)
	s41, _ := c4.Collect(*u41)

	assert.Equal(t, s41["active"], 1)
	assert.Equal(t, s41["accepts"], 8)
	assert.Equal(t, s41["handled"], 7)
	assert.Equal(t, s41["dropped"], 1)
	assert.Equal(t, s41["requests"], 23)
	assert.Equal(t, s41["current"], 4)
	assert.Equal(t, s41["reading"], 0)
	assert.Equal(t, s41["writing"], 1)
	assert.Equal(t, s41["waiting"], 2)
}
