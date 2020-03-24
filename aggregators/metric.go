package aggregators

import "github.com/mikanikos/WebsiteMonitoringTool/common"

// interface for metric aggregator
type MetricAggregator interface {
	AddMeasure(measure *common.WebsiteMeasure)
	RemoveMeasure(measure *common.WebsiteMeasure)
	GetMetric() string
}
