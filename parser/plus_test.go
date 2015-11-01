package parser

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestPlusParser(t *testing.T) {
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
			}
		}`
		fmt.Fprintln(w, payload)
	}))
	defer ts1.Close()

	p1 := &PlusParser{}
	u1, _ := url.Parse(ts1.URL)
	s1, _ := p1.Parse(u1.String())

	assert.Equal(t, 6, s1["version"])
	assert.Equal(t, 92676, s1["pid"])
	assert.Equal(t, map[string]interface{}{
		"accepted": 9202615,
		"dropped":  0,
		"active":   2,
		"idle":     17	,
	}, s1["connections"])
	assert.Equal(t, map[string]interface{}{
		"total":   19310289,
		"current": 2,
	}, s1["requests"])
	assert.NotNil(t, s1["server_zones"])
	assert.NotNil(t, s1["upstreams"])
	assert.NotNil(t, s1["caches"])
}
