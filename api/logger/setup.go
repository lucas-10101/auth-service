package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lucas-10101/auth-service/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Configure application logger and fallback based on properties
func Setup() {
	// configure fallback first
	api.FallbackLogger = getFallBackFileLogger()

	switch api.ApplicationProperties.LoggerProperties.LogDriver {
	case "file":
		api.Logger = getFileLogger()
	case "mongodb":
		api.Logger = getMongoDbLogger()
	case "stdout":
		api.Logger = getStdoutLogger()
	case "dummy":
		fallthrough
	default:
		api.Logger = getDummyLogger()
	}

}

// Default file log writter
func getFileLogger() *slog.Logger {
	file, err := os.OpenFile(api.ApplicationProperties.LoggerProperties.FileName, (os.O_APPEND | os.O_WRONLY | os.O_CREATE), 0640)

	if err != nil {
		panic(fmt.Sprintf("cant open file: %s (%s)", api.ApplicationProperties.LoggerProperties.FileName, err.Error()))
	}

	textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.Level(api.ApplicationProperties.LoggerProperties.LogLevel),
	})

	return slog.New(textHandler)
}

// MongoDB collection log writter
//
// if fallback is not enabled in properties, raise panic on connection error
func getMongoDbLogger() *slog.Logger {
	clientOptions := options.Client().
		ApplyURI(api.ApplicationProperties.LoggerProperties.MongoDbConnectionString).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMaxPoolSize(api.ApplicationProperties.LoggerProperties.MaxPoolSize).
		SetMinPoolSize(api.ApplicationProperties.LoggerProperties.MinPoolSize).
		SetAppName(api.ApplicationProperties.AppName)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		if api.ApplicationProperties.LoggerProperties.AllowFallback {
			return api.FallbackLogger
		}
		panic(fmt.Sprintf("cant create connection to mongodb log server (%s)", err.Error()))
	}

	err = client.Ping(context.Background(), readpref.Nearest())
	if err != nil {
		if api.ApplicationProperties.LoggerProperties.AllowFallback {
			return api.FallbackLogger
		}
		panic(fmt.Sprintf("cant reach mongodb log server (%s)", err.Error()))
	}

	logHandler := &MongoDBLogHandler{
		client: client,
		level:  slog.Level(api.ApplicationProperties.LoggerProperties.LogLevel),
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
func getFallBackFileLogger() *slog.Logger {

	if !api.ApplicationProperties.LoggerProperties.AllowFallback {
		return getDummyLogger()
	}

	file, err := os.OpenFile("fallback.log", (os.O_APPEND | os.O_WRONLY | os.O_CREATE), 0640)

	if err != nil {
		return getStdoutLogger()
	}

	textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.Level(api.ApplicationProperties.LoggerProperties.LogLevel),
	})

	return slog.New(textHandler)
}
