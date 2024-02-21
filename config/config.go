package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Core struct {
		Cron     string `yaml:"cron"`
		Location string `yaml:"location"`
		Logging  struct {
			// debug, info, warn, error
			Level *string `yaml:"level"`
			// console, json
			Type *string `yaml:"type"`
			// unix, rfc3339
			TimeFormat *string `yaml:"time_format"`
		} `yaml:"logging"`
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
	var logLevel zerolog.Level
	var logTimeFormat string

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

	switch {
	case config.Core.Logging.Level == nil:
		logLevel = zerolog.InfoLevel
	case *config.Core.Logging.Level == "debug":
		logLevel = zerolog.DebugLevel
	case *config.Core.Logging.Level == "info":
		logLevel = zerolog.InfoLevel
	case *config.Core.Logging.Level == "warn":
		logLevel = zerolog.WarnLevel
	case *config.Core.Logging.Level == "error":
		logLevel = zerolog.ErrorLevel
	default:
		return nil, fmt.Errorf("invalid log level: %s", *config.Core.Logging.Level)
	}

	switch {
	case config.Core.Logging.TimeFormat == nil:
		logTimeFormat = zerolog.TimeFormatUnix
	case *config.Core.Logging.TimeFormat == "unix":
		logTimeFormat = zerolog.TimeFormatUnix
	case *config.Core.Logging.TimeFormat == "rfc3339":
		logTimeFormat = time.RFC3339
	default:
		return nil, fmt.Errorf("invalid log time format: %s", *config.Core.Logging.TimeFormat)
	}

	switch {
	case config.Core.Logging.Type == nil:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logTimeFormat}).Level(logLevel)
	case *config.Core.Logging.Type == "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logTimeFormat}).Level(logLevel)
	case *config.Core.Logging.Type == "json":
		log.Logger = zerolog.New(os.Stderr).Level(logLevel)
	default:
		return nil, fmt.Errorf("invalid log type: %s", *config.Core.Logging.Type)
	}

	return config, nil
}
