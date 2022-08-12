package serviceDiscovery

import (
	"log"
	"simpleTool/simpleLogger"
)

type ServerDiscoveryLogger struct {
	logger simpleLogger.Logger
}

func NewServerDiscoveryLogger() *ServerDiscoveryLogger {
	return &ServerDiscoveryLogger{
		logger: simpleLogger.NewLogger(
			simpleLogger.WithCallPath(5),
			simpleLogger.WithFlag(log.Lshortfile|log.LstdFlags),
		),
	}
}

func (s *ServerDiscoveryLogger) Printf(layout string, v ...interface{}) {
	s.logger.Logf(layout, v...)
}
