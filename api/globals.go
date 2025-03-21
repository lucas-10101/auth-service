package api

import (
	"log/slog"

	"github.com/lucas-10101/auth-service/api/conf"
)

var (
	ApplicationProperties *conf.Properties = &conf.Properties{}
	Logger                *slog.Logger
	FallbackLogger        *slog.Logger
)
