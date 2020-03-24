package main

import (
	"flag"
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"github.com/mikanikos/WebsiteMonitoringTool/manager"
	"github.com/mikanikos/WebsiteMonitoringTool/statsManager"
	"github.com/mikanikos/WebsiteMonitoringTool/utils"
	"os"
)

func main() {

	configFile := flag.String("config", "config.yaml", "configuration file with the user parameters")

	// parse arguments
	flag.Parse()

	// parse config
	config, err := utils.ParseConfigFile(*configFile)
	if err != nil {
		fmt.Printf("exiting config file not parsed correctly: %s", err)
		os.Exit(1)
	}

	// create cui skeleton
	ui := cui.NewCui(config.StatsConfigs)

	// start monitoring and processing measures
	managers := make([]*manager.WebsiteManager, 0)
	for _, wConf := range config.WebsitesConfigs {
		wm := manager.NewWebsiteManager(wConf, config.AlertConfig, ui)
		wm.Run(config.StatsConfigs)
		managers = append(managers, wm)
	}

	// creates new stats manager
	sm := statsManager.NewStatsManager(config.StatsConfigs, managers, ui)
	// run the stats manger to periodically retrieve stats
	sm.Run()

	// start cui
	ui.StartCui()
}
