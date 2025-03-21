package conf

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

type LoggerProperties struct {
	// General
	LogLevel      int
	LogDriver     string
	AllowFallback bool // If log fails or is unavailable, use the following loggers in order: file, stdout. False is dummylogger

	// File specific
	FileName string

	// MongoDB Specific
	DatabaseName            string
	CollectionName          string
	MongoDbConnectionString string
	MaxPoolSize             uint64
	MinPoolSize             uint64
}
