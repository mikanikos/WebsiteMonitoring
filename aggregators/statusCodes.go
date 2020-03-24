package aggregators

import (
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
)

// status code aggregator computes stats for status codes response count
type StatusCodesAggregator struct {
	statusCodes map[int]uint64
}

func NewStatusCodesAggregator() *StatusCodesAggregator {
	return &StatusCodesAggregator{
		statusCodes: make(map[int]uint64),
	}
}

func (aggregator *StatusCodesAggregator) AddMeasure(measure *common.WebsiteMeasure) {
	aggregator.statusCodes[measure.StatusCode]++
}

func (aggregator *StatusCodesAggregator) RemoveMeasure(measure *common.WebsiteMeasure) {
	aggregator.statusCodes[measure.StatusCode]--
}

// get metric as a string (leaving map is intended)
func (aggregator *StatusCodesAggregator) GetMetric() string {
	if len(aggregator.statusCodes) != 0 {
		return fmt.Sprint(aggregator.statusCodes)
	}
	return "N/A"
}
