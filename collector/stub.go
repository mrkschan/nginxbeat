package collector

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/elastic/beats/libbeat/logp"
)

// StubCollector is a Collector that collects Nginx stub status page.
type StubCollector struct {
	http     *http.Client
	requests int
}

// NewStubCollector constructs a new StubCollector.
func NewStubCollector() Collector {
	return &StubCollector{
		http:     HTTPClient(),
		requests: 0,
	}
}

// Collect Nginx stub status from given url.
func (c *StubCollector) Collect(u url.URL) (map[string]interface{}, error) {
	res, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	// Nginx stub status sample:
	// Active connections: 1
	// server accepts handled requests
	//  7 7 19
	// Reading: 0 Writing: 1 Waiting: 0
	var re *regexp.Regexp
	scanner := bufio.NewScanner(res.Body)

	// Parse active connections.
	scanner.Scan()
	re = regexp.MustCompile("Active connections: (\\d+)")
	var active int
	if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
		logp.Warn("Fail to parse active connections from Nginx stub status")
		active = -1
	} else {
		active, _ = strconv.Atoi(matches[1])
	}

	// Skip request status headers.
	scanner.Scan()

	// Parse request status.
	scanner.Scan()
	re = regexp.MustCompile("\\s(\\d+)\\s+(\\d+)\\s+(\\d+)")
	var (
		accepts  int
		handled  int
		dropped  int
		requests int
		current  int
	)
	if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
		logp.Warn("Fail to parse request status from Nginx stub status")
		accepts = -1
		handled = -1
		dropped = -1
		requests = -1
		current = -1
	} else {
		accepts, _ = strconv.Atoi(matches[1])
		handled, _ = strconv.Atoi(matches[2])
		requests, _ = strconv.Atoi(matches[3])

		dropped = accepts - handled
		current = requests - c.requests

		c.requests = requests
	}

	// Parse connection status.
	scanner.Scan()
	re = regexp.MustCompile("Reading: (\\d+) Writing: (\\d+) Waiting: (\\d+)")
	var (
		reading int
		writing int
		waiting int
	)
	if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
		logp.Warn("Fail to parse connection status from Nginx stub status")
		reading = -1
		writing = -1
		waiting = -1
	} else {
		reading, _ = strconv.Atoi(matches[1])
		writing, _ = strconv.Atoi(matches[2])
		waiting, _ = strconv.Atoi(matches[3])
	}

	return map[string]interface{}{
		"active":   active,
		"accepts":  accepts,
		"handled":  handled,
		"dropped":  dropped,
		"requests": requests,
		"current":  current,
		"reading":  reading,
		"writing":  writing,
		"waiting":  waiting,
	}, nil
}
