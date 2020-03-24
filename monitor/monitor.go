package monitor

import (
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"net/http"
	"time"
)

const maxChannelSize = 100

// Monitor periodically checks a website to get measures
type Monitor struct {
	WebsiteConfig *common.WebsiteConfig
	MeasuresChan  chan *common.WebsiteMeasure
}

// NewMonitor creates a new monitor for a website
func NewMonitor(config *common.WebsiteConfig) *Monitor {
	return &Monitor{
		WebsiteConfig: config,
		MeasuresChan:  make(chan *common.WebsiteMeasure, maxChannelSize),
	}
}

// periodically checks the website
func (monitor *Monitor) PeriodicallyMonitorWebsite() {

	// set max timeout to prevent endless waiting, should be smaller than checkInterval
	maxTimeout := time.Duration(monitor.WebsiteConfig.CheckInterval/2) * time.Second
	client := http.Client{
		Timeout: maxTimeout,
	}

	// set check interval
	ticker := time.NewTicker(time.Duration(monitor.WebsiteConfig.CheckInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			start := time.Now()
			response, err := client.Get(monitor.WebsiteConfig.Url)
			elapsedTime := time.Since(start)
			if err != nil {
				monitor.MeasuresChan <- common.NewWebsiteMeasure(monitor.WebsiteConfig.Url, -1, elapsedTime, start)
			} else {
				response.Body.Close()
				monitor.MeasuresChan <- common.NewWebsiteMeasure(monitor.WebsiteConfig.Url, response.StatusCode, elapsedTime, start)
			}
		}
	}
}
