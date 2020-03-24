package metrics

import (
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"time"
)

// response time aggregator computes aggregated stats for average response time
type ResponseTimeAggregator struct {
	sumResponseTime time.Duration
	countMeasures   uint64
}

func NewResponseTimeAggregator() *ResponseTimeAggregator {
	return &ResponseTimeAggregator{
		sumResponseTime: time.Duration(0),
		countMeasures:   0,
	}
}

func (aggregator *ResponseTimeAggregator) AddMeasure(measure *common.WebsiteMeasure) {
	aggregator.sumResponseTime += measure.ResponseTime
	aggregator.countMeasures++
}

func (aggregator *ResponseTimeAggregator) RemoveMeasure(measure *common.WebsiteMeasure) {
	aggregator.sumResponseTime -= measure.ResponseTime
	aggregator.countMeasures--
}

// get metric as a string
func (aggregator *ResponseTimeAggregator) GetMetric() string {
	if aggregator.countMeasures != 0 {
		return fmt.Sprint(time.Duration(aggregator.sumResponseTime.Seconds()/float64(aggregator.countMeasures)*1000) * time.Millisecond)
	}
	return "N/A"
}
