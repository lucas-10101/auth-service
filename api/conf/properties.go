package conf

import (
	"log/slog"
)

// General application properties

type Properties struct {
	AppName string

	ServerProperties ServerProperties
	LoggerProperties LoggerProperties
}

// Server properties and definitions

type ServerProperties struct {
	Address            string
	Port               int
	TlsKeyPath         string
	TlsCertificatePath string
}

// Log properties and definitions

type LogDriver string

const (
	MongoDBLogDriver = LogDriver("mongodb")
	FileLogDriver    = LogDriver("file")
	DummyLogDriver   = LogDriver("dummy")
)

type LoggerProperties struct {
	// General
	LogLevel  slog.Level
	LogDriver LogDriver

	// File specific
	FileName string

	// MongoDB Specific
	MongoDbConnectionString string
	MaxPoolSize             uint64
	MinPoolSize             uint64
}
