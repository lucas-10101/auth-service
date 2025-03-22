package utils

type ErrorMessage string

const (

	// Properties
	PROPERTIES_FILE_READ_ERROR  = ErrorMessage("cannot open properties file, if dont exists run with --make-properties")
	PROPERTIES_ENTRY_BAD_FORMAT = ErrorMessage("malformed entry on properties file")

	// Logging
	LOG_FILE_READ_ERROR                    = ErrorMessage("cannot open specified log file")
	MONGODB_LOG_SERVER_CONNECTION_ERROR    = ErrorMessage("cannot create connection with mongodb log server")
	MONGODB_LOG_SERVER_COMMUNICATION_ERROR = ErrorMessage("cannot create connection with mongodb log server")

	// Server
	HTTPS_SERVER_START_FAILURE            = ErrorMessage("cannot start https server")
	HTTP_SERVER_START_FAILURE             = ErrorMessage("cannot start http server")
	HTTPS_REDIRECT_CONFIGURATION_MISMATCH = ErrorMessage("https is disabled, cannot create redirect")
)
