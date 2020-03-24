package common

import "time"

// WebsiteMeasure contains measures of a website activity
type WebsiteMeasure struct {
	Url          string
	StatusCode   int
	ResponseTime time.Duration
	MeasureTime  time.Time
}

func NewWebsiteMeasure(url string, status int, respTime time.Duration, measureTime time.Time) *WebsiteMeasure {
	return &WebsiteMeasure{
		Url:          url,
		StatusCode:   status,
		ResponseTime: respTime,
		MeasureTime:  measureTime,
	}
}
