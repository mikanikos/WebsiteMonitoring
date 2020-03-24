package stats

import (
	"sync"
	"time"

	"github.com/mikanikos/WebsiteMonitoringTool/aggregators"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
)

// WebsiteStats is a new object that computes aggregated stats based on a user-defined specific
type WebsiteStats struct {
	website   string
	timeFrame time.Duration
	measures  []*common.WebsiteMeasure

	availability    *aggregators.AvailabilityAggregator
	responseTime    *aggregators.ResponseTimeAggregator
	maxResponseTime *aggregators.MaxResponseTimeAggregator
	statusCodes     *aggregators.StatusCodesAggregator

	mutex sync.RWMutex
}

// StatsHandler manages the data coming from the monitor and computes stats
func NewWebsiteStats(url string, timeFrame float64) *WebsiteStats {
	return &WebsiteStats{
		website:         url,
		timeFrame:       time.Duration(timeFrame) * time.Second,
		measures:        make([]*common.WebsiteMeasure, 0),
		availability:    aggregators.NewAvailabilityAggregator(),
		responseTime:    aggregators.NewResponseTimeAggregator(),
		maxResponseTime: aggregators.NewMaxResponseTimeAggregator(),
		statusCodes:     aggregators.NewStatusCodesAggregator(),
	}
}

func (stats *WebsiteStats) PurgeOutdatedMeasures() {

	stats.mutex.Lock()
	defer stats.mutex.Unlock()

	startWindowPointer := 0
	for _, measure := range stats.measures {
		if time.Since(measure.MeasureTime) > stats.timeFrame {
			startWindowPointer++
			stats.availability.RemoveMeasure(measure)
			stats.responseTime.RemoveMeasure(measure)
			stats.statusCodes.RemoveMeasure(measure)
			stats.maxResponseTime.RemoveMeasure(measure)
		} else {
			break
		}
	}

	stats.measures = stats.measures[startWindowPointer:]
}

// StatsHandler manages the data coming from the monitor and computes stats
func (stats *WebsiteStats) UpdateStats(measure *common.WebsiteMeasure) {

	stats.mutex.Lock()
	defer stats.mutex.Unlock()

	//if strings.Contains(stats.website, "amazon") {
	//	fmt.Println("Before updating: ")
	//	fmt.Println(len(stats.measures))
	//}
	stats.measures = append(stats.measures, measure)
	//if strings.Contains(stats.website, "amazon") {
	//	fmt.Println("After updating: ")
	//	fmt.Println(len(stats.measures))
	//}

	stats.availability.AddMeasure(measure)
	stats.responseTime.AddMeasure(measure)
	stats.statusCodes.AddMeasure(measure)
	stats.maxResponseTime.AddMeasure(measure)

	//if strings.Contains(stats.website, "amazon") {
	//	fmt.Println("After method updating: ")
	//	fmt.Println(len(stats.measures))
	//}
}

func (stats *WebsiteStats) GetStatsObject() *cui.StatsObject {

	stats.PurgeOutdatedMeasures()

	statsObj := &cui.StatsObject{
		Website:         stats.website,
		Availability:    stats.availability.GetMetric(),
		AvgResponseTime: stats.responseTime.GetMetric(),
		MaxResponseTime: stats.maxResponseTime.GetMetric(),
		StatusCodes:     stats.statusCodes.GetMetric(),
	}

	return statsObj
}
