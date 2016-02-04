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

func TestPlusCollector(t *testing.T) {
	// It should report stats.
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := `{
			"version":6,
			"nginx_version":"1.9.4",
			"address":"206.251.255.64",
			"generation":12,
			"load_timestamp":1446285600278,
			"timestamp":1446296737353,
			"pid":92676,
			"processes":{"respawned":0},
			"connections":{"accepted":9202615,"dropped":0,"active":2,"idle":17},
			"ssl":{"handshakes":1165,"handshakes_failed":88,"session_reuses":292},
			"requests":{"total":19310289,"current":2},
			"server_zones":{
				"hg.nginx.org":{
					"processing":0,"requests":4182,
					"responses":{"1xx":0,"2xx":3760,"3xx":134,"4xx":68,"5xx":13,"total":3975},
					"discarded":207,"received":1159866,"sent":124423703
				}
			},
			"upstreams":{
				"hg-backend":{"peers":[
					{"id":0,"server":"10.0.0.1:8088","backup":false,"weight":5,"state":"up","active":0,"requests":3732,
					 "responses":{"1xx":0,"2xx":3662,"3xx":0,"4xx":68,"5xx":2,"total":3732},
					 "sent":1066721,"received":131453686,"fails":0,"unavail":0,
					 "health_checks":{"checks":1111,"fails":0,"unhealthy":0,"last_passed":true},
					 "downtime":0,"downstart":0,"selected":1446296711000
					 },
					{"id":1,"server":"10.0.0.1:8089","backup":true,"weight":1,"state":"unhealthy","active":0,"requests":0,
					 "responses":{"1xx":0,"2xx":0,"3xx":0,"4xx":0,"5xx":0,"total":0},
					 "sent":0,"received":0,"fails":0,"unavail":0,
					 "health_checks":{"checks":1114,"fails":1114,"unhealthy":1,"last_passed":false},
					 "downtime":11136335,"downstart":1446285601018,"selected":0
					 }
				],"keepalive":0}
			},
			"caches":{
				"http_cache":{
					"size":536649728,"max_size":536870912,"cold":false,
					"hit":{"responses":214171,"bytes":19170227606},
					"stale":{"responses":0,"bytes":0},
					"updating":{"responses":0,"bytes":0},
					"revalidated":{"responses":0,"bytes":0},
					"miss":{"responses":518012,"bytes":18427566829,"responses_written":316905,"bytes_written":8994873275},
					"expired":{"responses":41769,"bytes":2293640551,"responses_written":40536,"bytes_written":2268891511},
					"bypass":{"responses":55693,"bytes":2524371718,"responses_written":55677,"bytes_written":2524367270}
				}
			},
			"stream":{"server_zones":{"postgresql_loadbalancer":{"processing":0,"connections":74117,"received":7782285,"sent":418250173}},
					  "upstreams":{"postgresql_backends":{"peers":[
					  	{"id":0,"server":"10.0.0.2:15432","backup":false,"weight":1,"state":"up",
						 "active":0,"max_conns":42,"connections":24706,"connect_time":1,
						 "first_byte_time":1,"response_time":1,"sent":2594130,"received":139418346,
						 "fails":0,"unavail":0,"health_checks":{"checks":14892,"fails":0,"unhealthy":0,"last_passed":true},
						 "downtime":0,"downstart":0,"selected":1446360135000
						},
						{"id":1,"server":"10.0.0.2:15433","backup":false,"weight":1,"state":"up",
						 "active":0,"connections":24706,"connect_time":1,
						 "first_byte_time":1,"response_time":1,"sent":2594130,"received":139418665,
						 "fails":0,"unavail":0,"health_checks":{"checks":14892,"fails":0,"unhealthy":0,"last_passed":true},
						 "downtime":0,"downstart":0,"selected":1446360136000
						},
						{"id":2,"server":"10.0.0.2:15434","backup":false,"weight":1,"state":"up",
						 "active":0,"connections":24705,"connect_time":1,
						 "first_byte_time":1,"response_time":1,"sent":2594025,"received":139413162,
						 "fails":0,"unavail":0,"health_checks":{"checks":14892,"fails":0,"unhealthy":0,"last_passed":true},
						 "downtime":0,"downstart":0,"selected":1446360134000
						}
					  ]}}
			}
		}`
		fmt.Fprintln(w, payload)
	}))
	defer ts1.Close()

	c1 := &PlusCollector{}
	u1, _ := url.Parse(ts1.URL)
	s1, _ := c1.Collect(*u1)

	assert.Equal(t, 6, s1["version"])
	assert.Equal(t, 92676, s1["pid"])
	assert.Equal(t, map[string]interface{}{
		"accepted": 9202615,
		"dropped":  0,
		"active":   2,
		"idle":     17,
	}, s1["connections"])
	assert.Equal(t, map[string]interface{}{
		"total":   19310289,
		"current": 2,
	}, s1["requests"])
	assert.NotNil(t, s1["server_zones"])
	assert.IsType(t, []interface{}{}, s1["server_zones"])
	assert.NotNil(t, s1["upstreams"])
	assert.IsType(t, []interface{}{}, s1["upstreams"])
	assert.NotNil(t, s1["caches"])
	assert.IsType(t, []interface{}{}, s1["caches"])
	assert.NotNil(t, s1["stream"])
	stream1 := s1["stream"].(map[string]interface{})
	assert.IsType(t, []interface{}{}, stream1["server_zones"])
	assert.IsType(t, []interface{}{}, stream1["upstreams"])

	// It should report unexpected status code.
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failed", http.StatusInternalServerError)
	}))
	defer ts2.Close()

	c2 := &PlusCollector{}
	u2, _ := url.Parse(ts2.URL)
	s2, e2 := c2.Collect(*u2)

	assert.Nil(t, s2)
	assert.EqualError(t, e2, "HTTP500 Internal Server Error")

	// It should report malformed Nginx response.
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "malformed-json")
	}))
	defer ts3.Close()

	c3 := &PlusCollector{}
	u3, _ := url.Parse(ts3.URL)
	s3, e3 := c3.Collect(*u3)

	assert.Nil(t, s3)
	assert.Error(t, e3)
}
