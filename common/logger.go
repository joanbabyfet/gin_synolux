package common

import "go.uber.org/zap"

var Log *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	Log = logger.Sugar()
}