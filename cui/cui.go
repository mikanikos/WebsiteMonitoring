package cui

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
)

const (
	alertId        = "Alerts"
	defaultMessage = "No data available"
)

// Cui manages the ui of the application
type Cui struct {
	ui           *gocui.Gui
	statsConfigs []*common.StatConfig
}

// Cui creates a new gui
func NewCui(statsConfigs []*common.StatConfig) *Cui {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	ui := &Cui{
		ui:           gui,
		statsConfigs: statsConfigs,
	}

	gui.SetManagerFunc(ui.layout)

	if err := ui.SetKeyBindings(); err != nil {
		log.Fatalln(err)
	}
	return ui
}

// StartCui stats a new ui
func (cui *Cui) StartCui() {
	defer cui.ui.Close()
	if err := cui.ui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// UpdateStatsView updates the stats on a specific view
func (cui *Cui) UpdateStatsView(statsConfig *common.StatConfig, statsList []*StatsObject) {
	cui.ui.Update(func(g *gocui.Gui) error {
		view, err := cui.ui.View(statsConfig.String())
		if err != nil {
			return err
		}

		view.Clear()

		fmt.Fprintln(view)
		// sort for the same order
		sort.Slice(statsList, func(i, j int) bool { return strings.Compare(statsList[i].Website, statsList[j].Website) < 0 })

		format := "%-50s %-30s %-30s %-30s %-30s\n"
		stringFormatted := fmt.Sprintf(format, "website", "availability", "avg RT", "max RT", "response codes count")
		fmt.Fprintf(view, stringFormatted)
		fmt.Fprintln(view)
		fmt.Fprintln(view)
		// print all the websites stats
		for _, stats := range statsList {
			line := fmt.Sprintf(format, stats.Website, fmt.Sprint(stats.Availability), fmt.Sprint(stats.AvgResponseTime), fmt.Sprint(stats.MaxResponseTime), fmt.Sprint(stats.StatusCodes))
			fmt.Fprintln(view, line)
		}

		return nil
	})
}

// UpdateAlertsView adds a new alert message
func (cui *Cui) UpdateAlertsView(message string) {
	cui.ui.Update(func(g *gocui.Gui) error {
		v, err := g.View(alertId)
		if err != nil {
			return err
		}

		if strings.ContainsAny(message, "up") {
			fmt.Fprintln(v, "\x1b[0;32m"+message)
		} else {
			fmt.Fprintln(v, "\x1b[0;31m"+message)
		}

		return nil
	})
}

// layout of the ui, depends on the number of stats configured by the user
func (cui *Cui) layout(g *gocui.Gui) error {
	dimX, dimY := g.Size()
	maxX, maxY := float64(dimX), float64(dimY)
	numStats := float64(len(cui.statsConfigs))
	for it, conf := range cui.statsConfigs {
		i := float64(it)
		if v, err := g.SetView(conf.String(), 0, int(maxY*(i/(numStats+1))), int(maxX-1), int((maxY*((i+1)/(numStats+1)))-1)); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = conf.String()
			fmt.Fprintln(v, defaultMessage)
		}
	}

	if v, err := g.SetView(alertId, 0, int(maxY*(numStats/(numStats+1))), int(maxX-1), int(maxY-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = alertId
	}
	return nil
}

// set keybindings, can use mouse or keys to scroll
func (cui *Cui) SetKeyBindings() error {
	cui.ui.Mouse = true
	if err := cui.ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := cui.ui.SetKeybinding("", gocui.MouseWheelUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := cui.ui.SetKeybinding("", gocui.MouseWheelDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}

	if err := cui.ui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := cui.ui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
