package aggregators

import (
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"sort"
	"time"
)

// max response aggregator computes aggregated stats for max response time
type MaxResponseTimeAggregator struct {
	sortedResponseTimes []time.Duration
}

func NewMaxResponseTimeAggregator() *MaxResponseTimeAggregator {
	return &MaxResponseTimeAggregator{
		sortedResponseTimes: make([]time.Duration, 0),
	}
}

func (aggregator *MaxResponseTimeAggregator) AddMeasure(measure *common.WebsiteMeasure) {
	// binary search to find position where to insert element in sorted array
	i := sort.Search(len(aggregator.sortedResponseTimes), func(i int) bool { return aggregator.sortedResponseTimes[i] >= measure.ResponseTime })
	aggregator.sortedResponseTimes = append(aggregator.sortedResponseTimes, 0)
	copy(aggregator.sortedResponseTimes[i+1:], aggregator.sortedResponseTimes[i:])
	aggregator.sortedResponseTimes[i] = measure.ResponseTime
}

func (aggregator *MaxResponseTimeAggregator) RemoveMeasure(measure *common.WebsiteMeasure) {
	// binary search to find element to remove
	i := sort.Search(len(aggregator.sortedResponseTimes), func(i int) bool { return aggregator.sortedResponseTimes[i] == measure.ResponseTime })
	if i < len(aggregator.sortedResponseTimes) {
		copy(aggregator.sortedResponseTimes[i:], aggregator.sortedResponseTimes[i+1:])
		aggregator.sortedResponseTimes[len(aggregator.sortedResponseTimes)-1] = 0
		aggregator.sortedResponseTimes = aggregator.sortedResponseTimes[:len(aggregator.sortedResponseTimes)-1]
	}
}

// get metric as a string
func (aggregator *MaxResponseTimeAggregator) GetMetric() string {
	if len(aggregator.sortedResponseTimes) != 0 {
		return fmt.Sprint(aggregator.sortedResponseTimes[len(aggregator.sortedResponseTimes)-1])
	}
	return "N/A"
}
