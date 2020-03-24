package common

import "time"

// StatConfig for configuration of the stats
type StatConfig struct {
	CheckInterval     float64 `yaml:"displayInterval"`
	TimeFrameDuration float64 `yaml:"timeFrameDuration"`
}

// Creates a string representation of a statsConfig, it is used as key for the ui view and as title for simplicity
func (config *StatConfig) String() string {
	return "Stats for the past " + (time.Duration(config.TimeFrameDuration) * time.Second).String()
}
