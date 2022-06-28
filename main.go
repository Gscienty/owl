package main

import (
	"context"
	"flag"
	"owl/server/config"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func main() {
	ctx := context.Background()
	logger := logrus.WithContext(ctx).WithField("func", "main")

	logger.Info("starting owl-server")

	configFile := flag.String("file", "config.yml", "configuration file")

	// 初始化配置
	if err := config.Init(ctx, *configFile); err != nil {
		logger.WithError(err).Error("starting owl-server failed")
		return
	}
}
