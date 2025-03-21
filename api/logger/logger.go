package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lucas-10101/auth-service/api"
	"github.com/lucas-10101/auth-service/api/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetSelectedLogHandler() *slog.Logger {
	switch api.ApplicationProperties.LoggerProperties.LogDriver {
	case conf.FileLogDriver:
		return getFileLogHandler()
	case conf.MongoDBLogDriver:
		return getMongoDbLogHandler()
	case conf.DummyLogDriver:
		return getDummyLogHandler()
	default:
		return nil
	}
}

func getFileLogHandler() *slog.Logger {
	file, err := os.OpenFile(api.ApplicationProperties.LoggerProperties.FileName, (os.O_APPEND | os.O_WRONLY | os.O_CREATE), 640)

	if err != nil {
		panic(fmt.Sprintf("cant open file: %s (%s)", api.ApplicationProperties.LoggerProperties.FileName, err.Error()))
	}

	textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: api.ApplicationProperties.LoggerProperties.LogLevel,
	})

	return slog.New(textHandler)
}

func getMongoDbLogHandler() *slog.Logger {
	clientOptions := options.Client().
		ApplyURI(api.ApplicationProperties.LoggerProperties.MongoDbConnectionString).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMaxPoolSize(api.ApplicationProperties.LoggerProperties.MaxPoolSize).
		SetMinPoolSize(api.ApplicationProperties.LoggerProperties.MinPoolSize).
		SetAppName(api.ApplicationProperties.AppName)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(fmt.Sprintf("cant create connection to mongodb log server (%s)", err.Error()))
	}

	err = client.Ping(context.Background(), readpref.Nearest())
	if err != nil {
		panic(fmt.Sprintf("cant reach mongodb log server (%s)", err.Error()))
	}

	logHandler := &MongoDBLogHandler{
		client: client,
		level:  api.ApplicationProperties.LoggerProperties.LogLevel,
	}

	return slog.New(logHandler)
}

func getDummyLogHandler() *slog.Logger {
	return slog.New(&DummyLogHandler{})
}
