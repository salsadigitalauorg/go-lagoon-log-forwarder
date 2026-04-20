package logger

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
)

var (
	addSource             bool
	applicationName       string
	hostname              string
	logChannel            string
	logHost               string
	logPort               int
	logType               string // should match namespace to create index 'application-logs-{logType}'
	messageVersion        int
	handlerType           HandlerType
	normalizationPatterns []string
	once                  sync.Once
)

// synchronizedUDPWriter ensures UDP writes happen serially
type synchronizedUDPWriter struct {
	conn io.WriteCloser
	mu   sync.Mutex
}

func (w *synchronizedUDPWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.Write(p)
}

func (w *synchronizedUDPWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.Close()
}

// Initialize creates a multiwriter logger (udp and stdout) and sets it as the default
// slog
func Initialize(cfg Config) error {

	hostname, _ = os.Hostname()
	messageVersion = 3

	if err := config(cfg); err != nil {
		slog.Error("configuration error", slog.Any("err", err))
		return fmt.Errorf("configuration error: %w", err)
	}

	once.Do(func() {
		var writer io.Writer = os.Stdout

		udpConnection, err := connect()
		if err != nil {
			slog.Warn("Failed to connect to UDP endpoint, logging to stdout only", "error", err)
		} else {
			// Wrap UDP connection with synchronized writer to ensure serial writes
			syncUDPWriter := &synchronizedUDPWriter{conn: udpConnection}
			writer = io.MultiWriter(os.Stdout, syncUDPWriter)
		}

		baseHandler := slog.NewJSONHandler(
			writer,
			&slog.HandlerOptions{
				AddSource:   addSource,
				Level:       slog.LevelDebug,
				ReplaceAttr: replaceAttr,
			},
		)

		// Wrap with appropriate handler based on configuration
		var handler slog.Handler = baseHandler
		if handlerType == HandlerTypeNormalizing && len(normalizationPatterns) > 0 {
			handler = NewNormalizingHandler(baseHandler, normalizationPatterns)
		}

		slogger := slog.New(handler).With(defaultAttrs()...)

		slog.SetDefault(slogger)
	})

	return nil
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
