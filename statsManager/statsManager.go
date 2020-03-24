package statsManager

import (
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"github.com/mikanikos/WebsiteMonitoringTool/manager"
	"time"
)

// StatsManager is responsible of retrieving stats for each website and update the views
type StatsManager struct {
	statsConfigs    []*common.StatConfig
	websiteManagers []*manager.WebsiteManager
	ui              *cui.Cui
}

// NewStatsManager creates anew StatsManager given the configs
func NewStatsManager(configs []*common.StatConfig, managers []*manager.WebsiteManager, ui *cui.Cui) *StatsManager {
	return &StatsManager{
		statsConfigs:    configs,
		websiteManagers: managers,
		ui:              ui,
	}
}

// Start periodically ask for stats from the WebsiteStatsHandlers
func (sm *StatsManager) Run() {
	for _, stat := range sm.statsConfigs {
		go sm.periodicallyDisplayStats(stat)
	}
}

func (sm *StatsManager) periodicallyDisplayStats(config *common.StatConfig) {
	// set check interval
	ticker := time.NewTicker(time.Duration(config.CheckInterval) * time.Second)

	for {
		select {
		case <-ticker.C:

			// get stats for each website
			webStatsList := make([]*cui.StatsObject, 0)
			for _, webManagers := range sm.websiteManagers {
				wStats := webManagers.WebsiteStatsMap[config.String()]
				webStatsList = append(webStatsList, wStats.GetStatsObject())
			}

			// update view
			sm.ui.UpdateStatsView(config, webStatsList)
		}
	}
}
