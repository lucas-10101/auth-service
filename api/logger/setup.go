package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lucas-10101/auth-service/api/conf"
	"github.com/lucas-10101/auth-service/api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	logger         *slog.Logger
	fallbackLogger *slog.Logger
)

// Configure application logger and fallback based on properties
func Setup() {
	// configure fallback first
	fallbackLogger = getFallbackLogger()

	switch conf.ApplicationProperties.LoggerProperties.LogDriver {
	case "file":
		logger = getFileLogger()
	case "mongodb":
		logger = getMongoDbLogger()
	case "stdout":
		logger = getStdoutLogger()
	case "dummy":
		fallthrough
	default:
		logger = getDummyLogger()
	}

}

// Default file log writter
func getFileLogger() *slog.Logger {
	file, err := os.OpenFile(conf.ApplicationProperties.LoggerProperties.FileName, (os.O_APPEND | os.O_WRONLY | os.O_CREATE), 0640)

	if err != nil {
		panic(utils.LOG_FILE_READ_ERROR.WithErrorMessage(err))
	}

	textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.Level(conf.ApplicationProperties.LoggerProperties.LogLevel),
	})

	return slog.New(textHandler)
}

// MongoDB collection log writter
//
// if fallback is not enabled in properties, raise panic on connection error
func getMongoDbLogger() *slog.Logger {
	clientOptions := options.Client().
		ApplyURI(conf.ApplicationProperties.LoggerProperties.MongoDbConnectionString).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMaxPoolSize(conf.ApplicationProperties.LoggerProperties.MongoDbMaxPoolSize).
		SetMinPoolSize(conf.ApplicationProperties.LoggerProperties.MongoDbMinPoolSize).
		SetAppName(conf.ApplicationProperties.AppName)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		if conf.ApplicationProperties.LoggerProperties.AllowFallback {
			return fallbackLogger
		}
		GetFallbackLogger().Error(utils.MONGODB_LOG_SERVER_CONNECTION_ERROR.WithErrorMessage(err))
	}

	err = client.Ping(context.Background(), readpref.Nearest())
	if err != nil {
		if conf.ApplicationProperties.LoggerProperties.AllowFallback {
			return fallbackLogger
		}
		GetFallbackLogger().Error(utils.MONGODB_LOG_SERVER_COMMUNICATION_ERROR.WithErrorMessage(err))
	}

	logHandler := &MongoDBLogHandler{
		client: client,
		level:  slog.Level(conf.ApplicationProperties.LoggerProperties.LogLevel),
	}

	return slog.New(logHandler)
}

// Do nothing logger
func getDummyLogger() *slog.Logger {
	return slog.New(&DummyLogHandler{})
}

// Use stdout as output file
func getStdoutLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// configure fallback log handler if property [allowfallback = true]
// otherwise uses dummy logger to do nothinh on fail
func getFallbackLogger() *slog.Logger {

	if !conf.ApplicationProperties.LoggerProperties.AllowFallback {
		return getDummyLogger()
	}

	file, err := os.OpenFile("fallback.log", (os.O_APPEND | os.O_WRONLY | os.O_CREATE), 0640)

	if err != nil {
		return getStdoutLogger()
	}

	textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.Level(conf.ApplicationProperties.LoggerProperties.LogLevel),
	})

	return slog.New(textHandler)
}

func GetLogger() *slog.Logger {
	return logger
}

func GetFallbackLogger() *slog.Logger {
	return fallbackLogger
}
