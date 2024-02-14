package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Core struct {
		Cron     string `yaml:"cron"`
		Location string `yaml:"location"`
	} `yaml:"core"`
	Notifications struct {
		Discord []struct {
			Name    string `yaml:"name"`
			Webhook string `yaml:"webhook"`
		} `yaml:"discord"`
	} `yaml:"notifications"`
	Repositories []struct {
		Name                 string   `yaml:"name"`
		Prerelease           bool     `yaml:"prerelease"`
		TagRegex             *string  `yaml:"tag_regex"`
		NotificationChannels []string `yaml:"notification_channels"`
	} `yaml:"repositories"`
}

func ParseConfig() (*Config, error) {
	config := &Config{}

	filename := "config.yml"
	if os.Getenv("CONFIG_FILE") != "" {
		filename = os.Getenv("CONFIG_FILE")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
