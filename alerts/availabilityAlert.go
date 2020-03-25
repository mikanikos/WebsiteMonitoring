package alerts

import (
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"github.com/mikanikos/WebsiteMonitoringTool/metrics"
)

// AvailabilityAlert is an alert for the availability metric
type AvailabilityAlert struct {
	website            string
	AlertConfig        *common.AlertConfig
	measures           []*common.WebsiteMeasure
	availability       *metrics.AvailabilityAggregator
	mutex              sync.RWMutex
	isWebsiteAvailable bool
	ui                 *cui.Cui
}

// NewAvailabilityAlert creates a new availability alert
func NewAvailabilityAlert(url string, config *common.AlertConfig, ui *cui.Cui) *AvailabilityAlert {
	return &AvailabilityAlert{
		website:            url,
		AlertConfig:        config,
		measures:           make([]*common.WebsiteMeasure, 0),
		availability:       metrics.NewAvailabilityAggregator(),
		isWebsiteAvailable: true,
		ui:                 ui,
	}
}

// PurgeOutdatedMeasures removes outdated data from current elements and updates the aggregators
func (alert *AvailabilityAlert) PurgeOutdatedMeasures() {

	alert.mutex.Lock()
	defer alert.mutex.Unlock()

	startWindowPointer := 0
	for _, measure := range alert.measures {
		if time.Since(measure.MeasureTime) > (time.Duration(alert.AlertConfig.TimeFrameDuration) * time.Second) {
			startWindowPointer++
			alert.availability.RemoveMeasure(measure)
		} else {
			// stop here because new elements are inserted at the tail
			break
		}
	}

	alert.measures = alert.measures[startWindowPointer:]
}

// UpdateStats adds a new measure and updates the aggregators
func (alert *AvailabilityAlert) UpdateStats(measure *common.WebsiteMeasure) {

	alert.mutex.Lock()
	defer alert.mutex.Unlock()

	alert.measures = append(alert.measures, measure)
	alert.availability.AddMeasure(measure)
}

// UpdateAlertData updates alert data after a new measure
func (alert *AvailabilityAlert) UpdateAlertData() {

	if !math.IsNaN(alert.availability.GetValue()) {
		isCurrentlyAvailable := alert.availability.GetValue() >= alert.AlertConfig.Threshold

		if (!isCurrentlyAvailable && alert.isWebsiteAvailable) || (isCurrentlyAvailable && !alert.isWebsiteAvailable) {
			alert.isWebsiteAvailable = !alert.isWebsiteAvailable
			alert.ui.UpdateAlertsView(alert.getAlertMessage(alert.isWebsiteAvailable))
		}
	}
}

// create message for alert
func (alert *AvailabilityAlert) getAlertMessage(status bool) string {
	statusString := "down"
	if status {
		statusString = "up"
	}

	var sb strings.Builder

	sb.WriteString("Website ")
	sb.WriteString(alert.website)
	sb.WriteString(" is ")
	sb.WriteString(statusString)
	sb.WriteString(". availability=")
	sb.WriteString(alert.availability.GetMetric())
	sb.WriteString(", time=")
	sb.WriteString(time.Now().Format(time.RFC3339))

	return sb.String()
}
