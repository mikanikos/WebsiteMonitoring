package stats

import (
	"sync"
	"time"

	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"github.com/mikanikos/WebsiteMonitoringTool/cui"
	"github.com/mikanikos/WebsiteMonitoringTool/metrics"
)

// WebsiteStatsHandler is a new object that computes aggregated stats for a website based on a user-defined specific
type WebsiteStatsHandler struct {
	website   string
	timeFrame time.Duration
	measures  []*common.WebsiteMeasure

	availability    *metrics.AvailabilityAggregator
	responseTime    *metrics.ResponseTimeAggregator
	maxResponseTime *metrics.MaxResponseTimeAggregator
	statusCodes     *metrics.StatusCodesAggregator

	mutex sync.RWMutex
}

// StatsHandler manages the data coming from the monitor and computes stats
func NewWebsiteStats(url string, timeFrame float64) *WebsiteStatsHandler {
	return &WebsiteStatsHandler{
		website:         url,
		timeFrame:       time.Duration(timeFrame) * time.Second,
		measures:        make([]*common.WebsiteMeasure, 0),
		availability:    metrics.NewAvailabilityAggregator(),
		responseTime:    metrics.NewResponseTimeAggregator(),
		maxResponseTime: metrics.NewMaxResponseTimeAggregator(),
		statusCodes:     metrics.NewStatusCodesAggregator(),
	}
}

func (stats *WebsiteStatsHandler) PurgeOutdatedMeasures() {

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
func (stats *WebsiteStatsHandler) UpdateStats(measure *common.WebsiteMeasure) {

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

func (stats *WebsiteStatsHandler) GetStatsObject() *cui.StatsObject {

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
