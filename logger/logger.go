package logger

import (
	"github.com/zhayt/user-service/config"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.AppMode == "dev" {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}
