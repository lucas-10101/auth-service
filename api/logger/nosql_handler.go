package logger

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBLogHandler struct {
	client *mongo.Client
	level  slog.Level
}

func (logger *MongoDBLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= logger.level
}

func (logger *MongoDBLogHandler) Handle(ctx context.Context, record slog.Record) error {

	logger.client.Ping(context.Background(), readpref.Nearest())

	return nil
}

func (logger *MongoDBLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return nil
}

func (logger *MongoDBLogHandler) WithGroup(name string) slog.Handler {
	return nil
}

func (logger *MongoDBLogHandler) getCollection() mongo.Collection {
	return mongo.Collection{}
}
