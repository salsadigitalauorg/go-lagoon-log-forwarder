package logger

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

// Initialize creates a multiwriter logger (udp and stdout) and sets it as the default
// slog
func Initialize() {

	validate()

	once.Do(func() {

		udpConnection, err := connect()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Failed to connect to UDP endpoint: %s", err)
			os.Exit(1)
		}

		slogger := slog.New(
			slog.NewJSONHandler(
				io.MultiWriter(
					os.Stdout,
					udpConnection,
				),
				&slog.HandlerOptions{
					AddSource:   addSource,
					Level:       slog.LevelDebug,
					ReplaceAttr: replaceAttr,
				},
			)).With(defaultAttrs()...)

		slog.SetDefault(slogger)
	})
}

func defaultAttrs() []any {

	return []any{
		slog.Int("@version", messageVersion),
		slog.String("application", applicationName),
		slog.String("channel", logChannel),
		slog.Group("context"),
		slog.Group("extra"),
		slog.String("host", hostname),
		// NOTE: Refactoring will be required if we want to override this per project
		slog.String("type", logType),
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 {
		switch a.Key {
		case "msg":
			a.Key = "message"
		case "time":
			a.Key = "@timestamp"
		case "timestampOverride":
			a.Key = "@timestamp"
		}
	}
	return a
}

func connect() (*net.UDPConn, error) {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", logHost, logPort))
	if err != nil {
		slog.Error("Failed to resolve udp address")
		return nil, err
	}

	con, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		slog.Error("Failed to dial udp")
		return nil, err
	}

	return con, nil

}
