package logging

import (
	"fmt"

	"github.com/cresta/zapctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogging(level string) (*zapctx.Logger, error) {
	cfg := zap.NewProductionConfig()
	if level != "" {
		var l zapcore.Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			return nil, fmt.Errorf("invalid logging level %s: %w", level, err)
		}
		cfg.Level.SetLevel(l)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("unable to setup logging: %w", err)
	}
	return zapctx.New(logger), nil
}
