package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/logp"
	"github.com/elastic/libbeat/publisher"
)

const selector = "nginxbeat"

// Nginxbeat implements Beater interface and sends Nginx status using libbeat.
type Nginxbeat struct {
	// NbConfig holds configurations of Nginxbeat parsed by libbeat.
	NbConfig ConfigSettings

	done     chan uint
	requests int
	events   publisher.Client

	url    *url.URL
	format string
	period time.Duration
}

// Config Nginxbeat according to nginxbeat.yml.
func (nb *Nginxbeat) Config(b *beat.Beat) error {
	err := cfgfile.Read(&nb.NbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	var u string
	if nb.NbConfig.Input.URL != "" {
		u = nb.NbConfig.Input.URL
	} else {
		u = "http://127.0.0.1/status"
	}
	nb.url, err = url.Parse(u)
	if err != nil {
		logp.Err("Invalid Nginx status page: %v", err)
		return err
	}

	var f string
	if nb.NbConfig.Input.Format != "" {
		f = nb.NbConfig.Input.Format
	} else {
		f = "stub"
	}
	if f != "stub" && f != "json" {
		err = fmt.Errorf("%s is an unsupported format", f)
	}
	if err != nil {
		logp.Err("Invalid Nginx status format: %v", err)
		return err
	}
	nb.format = f

	if nb.NbConfig.Input.Period != nil {
		nb.period = time.Duration(*nb.NbConfig.Input.Period) * time.Second
	} else {
		nb.period = 1 * time.Second
	}

	logp.Debug(selector, "Init nginxbeat")
	logp.Debug(selector, "Watch %v", nb.url)
	logp.Debug(selector, "Format %v", nb.format)
	logp.Debug(selector, "Period %v", nb.period)

	return nil
}

// Setup Nginxbeat.
func (nb *Nginxbeat) Setup(b *beat.Beat) error {
	nb.events = b.Events
	nb.requests = 0
	nb.done = make(chan uint)

	return nil
}

// Run Nginxbeat.
func (nb *Nginxbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run nginxbeat")

	ticker := time.NewTicker(nb.period)
	defer ticker.Stop()

	for {
		select {
		case <-nb.done:
			goto GotoFinish
		case <-ticker.C:
		}

		start := time.Now()

		if nb.format == "stub" {
			s, err := nb.getStubStatus()
			if err != nil {
				logp.Err("Fail to read Nginx stub status: %v", err)
				goto GotoNext
			}

			nb.events.PublishEvent(common.MapStr{
				"timestamp": common.Time(time.Now()),
				"type":      "nginx",
				"nginx":     s,
			})
		}

	GotoNext:
		end := time.Now()
		duration := end.Sub(start)
		if duration.Nanoseconds() > nb.period.Nanoseconds() {
			logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
		}
	}

GotoFinish:
	return nil
}

// Cleanup Nginxbeat.
func (nb *Nginxbeat) Cleanup(b *beat.Beat) error {
	return nil
}

// Stop Nginxbeat.
func (nb *Nginxbeat) Stop() {
	logp.Debug(selector, "Stop nginxbeat")
	close(nb.done)
}

func (nb *Nginxbeat) getStubStatus() (map[string]int, error) {
	res, err := http.Get(nb.url.String())
	if err != nil {
		return nil, err
	}

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
	defer res.Body.Close()

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
		current = requests - nb.requests

		nb.requests = requests
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

	return map[string]int{
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
