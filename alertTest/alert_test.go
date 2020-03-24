package alertTest

import (
	"context"
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"github.com/mikanikos/WebsiteMonitoringTool/manager"
	"github.com/mikanikos/WebsiteMonitoringTool/statsManager"
	"github.com/mikanikos/WebsiteMonitoringTool/utils"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAlert(t *testing.T) {

	// creates a simple routine to launch a server and make him available or not after some time
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		})

		server := http.Server{
			Addr: ":8080",
		}

		time.Sleep(15 * time.Second)
		go server.ListenAndServe()
		time.Sleep(30 * time.Second)
		server.Shutdown(context.Background())
		time.Sleep(25 * time.Second)
	}()

	// parse config
	config, err := utils.ParseConfigFile("config_test.yaml")
	if err != nil {
		fmt.Printf("exiting config file not parsed correctly: %s", err)
		os.Exit(1)
	}

	// create cui skeleton
	ui := cui.NewCui(config.StatsConfigs)

	// start monitoring
	managers := make([]*manager.WebsiteManager, 0)
	for _, wConf := range config.WebsitesConfigs {
		wm := manager.NewWebsiteManager(wConf, config.AlertConfig, ui)
		wm.Run(config.StatsConfigs)
		managers = append(managers, wm)
	}

	sm := statsManager.NewStatsManager(config.StatsConfigs, managers, ui)
	sm.Run()

	// start cui
	ui.StartCui()
}
