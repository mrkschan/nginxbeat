package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/logp"

	"github.com/mrkschan/nginxbeat/collector"
	"github.com/mrkschan/nginxbeat/publisher"
)

const selector = "nginxbeat"

// Nginxbeat implements Beater interface and sends Nginx status using libbeat.
type Nginxbeat struct {
	// NbConfig holds configurations of Nginxbeat parsed by libbeat.
	NbConfig ConfigSettings

	done chan uint

	urls   []*url.URL
	period time.Duration
}

// Config Nginxbeat according to nginxbeat.yml.
func (nb *Nginxbeat) Config(b *beat.Beat) error {
	err := cfgfile.Read(&nb.NbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	var urlConfig []string
	if nb.NbConfig.Input.URLs != nil {
		urlConfig = nb.NbConfig.Input.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1/status"}
	}

	nb.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid Nginx status page: %v", err)
			return err
		}
		switch u.Fragment {
		case "plus", "stub":
		default:
			err := fmt.Errorf("%s is not supported", u.Fragment)
			logp.Err("Invalid Nginx status page format: %v", err)
			return err
		}
		nb.urls[i] = u
	}

	if nb.NbConfig.Input.Period != nil {
		nb.period = time.Duration(*nb.NbConfig.Input.Period) * time.Second
	} else {
		nb.period = 10 * time.Second
	}

	logp.Debug(selector, "Init nginxbeat")
	logp.Debug(selector, "Watch %v", nb.urls)
	logp.Debug(selector, "Period %v", nb.period)

	return nil
}

// Setup Nginxbeat.
func (nb *Nginxbeat) Setup(b *beat.Beat) error {
	nb.done = make(chan uint)

	return nil
}

// Run Nginxbeat.
func (nb *Nginxbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run nginxbeat")

	for _, u := range nb.urls {
		go func(u *url.URL) {
			var c collector.Collector
			var p publisher.Publisher

			switch u.Fragment {
			case "stub":
				c = collector.NewStubCollector()
				p = publisher.NewStubPublisher(b.Events)
			case "plus":
				c = collector.NewPlusCollector()
				p = publisher.NewPlusPublisher(b.Events)
			}

			ticker := time.NewTicker(nb.period)
			defer ticker.Stop()

			for {
				select {
				case <-nb.done:
					goto GotoFinish
				case <-ticker.C:
				}

				start := time.Now()

				s, err := c.Collect(*u)
				if err != nil {
					logp.Err("Fail to read Nginx status: %v", err)
					goto GotoNext
				}
				p.Publish(s)

			GotoNext:
				end := time.Now()
				duration := end.Sub(start)
				if duration.Nanoseconds() > nb.period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}(u)
	}

	<-nb.done
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
