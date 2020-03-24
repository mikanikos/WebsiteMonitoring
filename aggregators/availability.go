package aggregators

import (
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"math"
)

// availability aggregator computes aggregated stats for availability
type AvailabilityAggregator struct {
	sumAvailability float64
	countMeasures   uint64
}

func NewAvailabilityAggregator() *AvailabilityAggregator {
	return &AvailabilityAggregator{
		sumAvailability: 0,
	}
}

func (aggregator *AvailabilityAggregator) AddMeasure(measure *common.WebsiteMeasure) {
	status := float64(0)
	if measure.StatusCode >= 200 && measure.StatusCode <= 299 {
		status = 1
	}
	aggregator.sumAvailability += status
	aggregator.countMeasures++
}

func (aggregator *AvailabilityAggregator) RemoveMeasure(measure *common.WebsiteMeasure) {
	status := float64(0)
	if measure.StatusCode >= 200 && measure.StatusCode <= 299 {
		status = 1
	}
	aggregator.sumAvailability -= status
	aggregator.countMeasures--
}

// get metric as a string
func (aggregator *AvailabilityAggregator) GetMetric() string {
	if aggregator.countMeasures != 0 {
		return fmt.Sprint(math.Round(aggregator.sumAvailability/float64(aggregator.countMeasures)*100)) + "%"
	}
	return "N/A"
}

// get metric as value
func (aggregator *AvailabilityAggregator) GetValue() float64 {
	if aggregator.countMeasures != 0 {
		return aggregator.sumAvailability / float64(aggregator.countMeasures) * 100
	}
	return math.NaN()
}
