package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/logp"
	"github.com/elastic/libbeat/publisher"

	"github.com/mrkschan/nginxbeat/parser"
)

const selector = "nginxbeat"

// Nginxbeat implements Beater interface and sends Nginx status using libbeat.
type Nginxbeat struct {
	// NbConfig holds configurations of Nginxbeat parsed by libbeat.
	NbConfig ConfigSettings

	done   chan uint
	events publisher.Client

	urls   []*url.URL
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
		nb.urls[i] = u
	}

	var f string
	if nb.NbConfig.Input.Format != "" {
		f = nb.NbConfig.Input.Format
	} else {
		f = "stub"
	}
	if f != "stub" && f != "plus" {
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
		nb.period = 10 * time.Second
	}

	logp.Debug(selector, "Init nginxbeat")
	logp.Debug(selector, "Watch %v", nb.urls)
	logp.Debug(selector, "Format %v", nb.format)
	logp.Debug(selector, "Period %v", nb.period)

	return nil
}

// Setup Nginxbeat.
func (nb *Nginxbeat) Setup(b *beat.Beat) error {
	nb.events = b.Events
	nb.done = make(chan uint)

	return nil
}

// Run Nginxbeat.
func (nb *Nginxbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run nginxbeat")

	for _, u := range nb.urls {
		go func() {
			var p parser.Parser
			switch nb.format {
			case "stub":
				p = parser.NewStubParser()
			case "plus":
				p = parser.NewPlusParser()
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

				s, err := p.Parse(*u)
				if err != nil {
					logp.Err("Fail to read Nginx status: %v", err)
					goto GotoNext
				}
				nb.events.PublishEvent(common.MapStr{
					"@timestamp": common.Time(time.Now()),
					"type":       "nginx",
					"nginx":      s,
				})

			GotoNext:
				end := time.Now()
				duration := end.Sub(start)
				if duration.Nanoseconds() > nb.period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}()
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
