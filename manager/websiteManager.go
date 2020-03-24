package manager

import (
	"github.com/mikanikos/WebsiteMonitoringTool/alerts"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"github.com/mikanikos/WebsiteMonitoringTool/monitor"
	"github.com/mikanikos/WebsiteMonitoringTool/stats"
)

// WebsiteManager manages the data coming from the monitor
type WebsiteManager struct {
	monitor         *monitor.Monitor
	WebsiteStatsMap map[string]*stats.WebsiteStats
	alert           *alerts.AvailabilityAlert
}

// WebsiteManager manages the data coming from the monitor and computes stats
func NewWebsiteManager(config *common.WebsiteConfig, alertConfig *common.AlertConfig, ui *cui.Cui) *WebsiteManager {
	return &WebsiteManager{
		monitor:         monitor.NewMonitor(config),
		WebsiteStatsMap: make(map[string]*stats.WebsiteStats, 0),
		alert:           alerts.NewAvailabilityAlert(config.Url, alertConfig, ui),
	}
}

// Start allows the statshandler to start monitoring and processing data
func (wm *WebsiteManager) Run(statsConfigs []*common.StatConfig) {

	for _, sConf := range statsConfigs {
		wStats := stats.NewWebsiteStats(wm.monitor.WebsiteConfig.Url, sConf.TimeFrameDuration)
		wm.WebsiteStatsMap[sConf.String()] = wStats
	}

	go wm.processWebsiteMeasures()
	go wm.monitor.PeriodicallyMonitorWebsite()
}

// process new measure
func (wm *WebsiteManager) processWebsiteMeasures() {
	for newMeasure := range wm.monitor.MeasuresChan {

		for _, stat := range wm.WebsiteStatsMap {
			stat.PurgeOutdatedMeasures()
			stat.UpdateStats(newMeasure)
		}

		wm.alert.PurgeOutdatedMeasures()
		wm.alert.UpdateStats(newMeasure)
		wm.alert.UpdateAlertData()
	}
}
