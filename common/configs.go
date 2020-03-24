package common

// Config contains all the configuration parameters that are set by the user through a config file
type Config struct {
	WebsitesConfigs []*WebsiteConfig `yaml:"websitesConfigs"`
	StatsConfigs    []*StatConfig    `yaml:"statsConfigs"`
	AlertConfig     *AlertConfig     `yaml:"alertConfig"`
}

func NewConfig() *Config {
	return &Config{
		StatsConfigs:    make([]*StatConfig, 0),
		WebsitesConfigs: make([]*WebsiteConfig, 0),
	}
}
