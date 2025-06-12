package logger

import (
	"errors"
	"log/slog"
)

type Config struct {
	AddSource       bool
	ApplicationName string
	LogChannel      string
	LogHost         string
	LogPort         int
	LogType         string
	MessageVersion  int
}

// NewConfig returns a Config struct with default values
func NewConfig() Config {
	return Config{
		AddSource:       true,
		ApplicationName: "",
		LogChannel:      "LagoonLogs",
		LogHost:         "", // Will default to localhost in validation
		LogPort:         5140,
		LogType:         "", // Required - must be set by user
		MessageVersion:  1,
	}
}

func config(cfg Config) error {
	addSource = cfg.AddSource
	applicationName = cfg.ApplicationName
	logChannel = cfg.LogChannel
	logHost = cfg.LogHost
	logPort = cfg.LogPort
	logType = cfg.LogType
	messageVersion = cfg.MessageVersion
	return validate()
}

func validate() error {

	// validate logstashHost
	if len(logHost) == 0 {
		slog.Warn(
			"log.host is not supplied and will default to localhost",
		)
	}

	if len(logType) == 0 {
		return errors.New("logType is required")
	}

	return nil
}
