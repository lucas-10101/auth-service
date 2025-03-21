package logger

import (
	"context"
	"log/slog"

	"github.com/lucas-10101/auth-service/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// slog.Handler implementations to log in mongodb
type MongoDBLogHandler struct {
	client     *mongo.Client
	level      slog.Level
	groupName  string
	attributes []primitive.E
}

func (logger *MongoDBLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= logger.level
}

func (logger *MongoDBLogHandler) Handle(ctx context.Context, record slog.Record) error {
	collection := logger.getCollection()

	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "app-name", Value: api.ApplicationProperties.AppName},
		{Key: "time", Value: record.Time},
		{Key: "level", Value: record.Level},
		{Key: "group", Value: logger.groupName},
		{Key: "attributes", Value: bson.D(logger.attributes)},
		{Key: "message", Value: record.Message},
	})
	return err
}

func (logger *MongoDBLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {

	attributes := []primitive.E{}
	for _, e := range attrs {
		attributes = append(attributes, primitive.E{
			Key:   e.Key,
			Value: e.Value,
		})
	}

	return &MongoDBLogHandler{
		client:     logger.client,
		level:      logger.level,
		groupName:  logger.groupName,
		attributes: attributes,
	}
}

func (logger *MongoDBLogHandler) WithGroup(name string) slog.Handler {
	return &MongoDBLogHandler{
		client:     logger.client,
		level:      logger.level,
		groupName:  name,
		attributes: logger.attributes,
	}
}

func (logger *MongoDBLogHandler) getCollection() *mongo.Collection {
	return logger.client.
		Database(api.ApplicationProperties.LoggerProperties.DatabaseName).
		Collection(api.ApplicationProperties.LoggerProperties.CollectionName)
}
