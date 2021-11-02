package logging

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)
type contextKey = string

const loggerKey = contextKey("logger")


var (
	level             = zapcore.Level(-1)
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)


func NewLogger(level zapcore.Level) *zap.SugaredLogger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg := zap.Config{
		Encoding:         "console",
		EncoderConfig:    ec,
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	return logger.Sugar()
}

func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger(level)
	})
	return defaultLogger
}

func WithLogger(ctx context.Context,logger *zap.SugaredLogger) context.Context {
	if gCtx , ok := ctx.(*gin.Context); ok {
		ctx = gCtx.Request.Context()
	}
	return context.WithValue(ctx,loggerKey,logger)
}

func fromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil{
		return DefaultLogger()
	}
	if gCtx,ok := ctx.(*gin.Context);ok{
		ctx = gCtx.Request.Context()
	}
	if logger , ok := ctx.Value(loggerKey).(*zap.SugaredLogger);ok{
		return logger
	}
	return DefaultLogger()
}


