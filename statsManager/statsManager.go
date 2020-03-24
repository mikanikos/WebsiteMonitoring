package statsManager

import (
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"github.com/mikanikos/WebsiteMonitoringTool/manager"
	"time"
)

// StatsManager manages the data coming from the monitor
type StatsManager struct {
	statsConfigs    []*common.StatConfig
	websiteManagers []*manager.WebsiteManager
	ui              *cui.Cui
}

// WebsiteManager manages the data coming from the monitor and computes stats
func NewStatsManager(configs []*common.StatConfig, managers []*manager.WebsiteManager, ui *cui.Cui) *StatsManager {
	return &StatsManager{
		statsConfigs:    configs,
		websiteManagers: managers,
		ui:              ui,
	}
}

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

			webStatsList := make([]*cui.StatsObject, 0)
			for _, webManagers := range sm.websiteManagers {
				wStats := webManagers.WebsiteStatsMap[config.String()]
				webStatsList = append(webStatsList, wStats.GetStatsObject())
			}

			sm.ui.UpdateStatsView(config, webStatsList)
		}
	}
}
