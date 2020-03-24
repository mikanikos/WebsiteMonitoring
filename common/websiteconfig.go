package common

// WebsiteConfig for configuration of the website monitoring
type WebsiteConfig struct {
	Url           string  `yaml:"url"`
	CheckInterval float64 `yaml:"checkInterval"`
}
