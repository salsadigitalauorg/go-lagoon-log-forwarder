package logger

import (
	"flag"
	"log/slog"
	"os"
	"sync"
)

var (
	addSource       bool
	applicationName string
	hostname        string
	logChannel      string
	logHost         string
	logPort         int
	logType         string // should match namespace to create index 'application-logs-{logType}'
	messageVersion  int
	once            sync.Once
)

func init() {

	hostname, _ = os.Hostname()

	messageVersion = 3

	// Required
	flag.StringVar(&logType, "log.fields.type", "", "Log Type (must match existing k8s namespace)")
	flag.StringVar(&logHost, "log.host", "", "UDP host")
	// Optionals
	flag.BoolVar(&addSource, "log.addSource", true, "Add source to logs")
	flag.StringVar(&logChannel, "log.channel", "LagoonLogs", "Channel name")
	flag.StringVar(&applicationName, "log.fields.applicationName", "", "Application name")
	flag.IntVar(&logPort, "log.port", 5140, "UDP port")

}

func validate() {

	// validate logstashHost
	if len(logHost) == 0 {
		slog.Warn(
			"log.host is not supplied and will default to localhost",
		)
	}

	if len(logType) == 0 {
		slog.Error(
			"log.fields.type is required",
		)
		os.Exit(1)
	}

}
