package common

// AlertConfig for configuration of the alert
type AlertConfig struct {
	Threshold         float64 `yaml:"threshold"`
	TimeFrameDuration float64 `yaml:"timeFrameDuration"`
}
