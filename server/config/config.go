// 应用服务基础配置
package config

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

// AppConfig 应用配置
type AppConfig struct {
	Port uint16 `default:"1080"`

	WebSocket struct {
		HandshakeTimeoutMills int `default:"5000"`
		ReadBuffer            int `default:"1024"`
		WriteBuffer           int `default:"1024"`
		Subprotocols          []string
		CheckOrigin           bool `default:"false"`
		EnableCompression     bool `default:"false"`
	}
}

// appConfig 全局应用配置
var appConfig AppConfig

// Web Socket 相关配置
var websocketUpgrader websocket.Upgrader

// GetAppConfig 获取全局配置
func GetAppConfig(ctx context.Context) *AppConfig { return &appConfig }

// GetWebSocketUpgrader 获取 WebSocket Upgrader 配置
func GetWebSocketUpgrader(ctx context.Context) *websocket.Upgrader { return &websocketUpgrader }

// Init 初始化配置
//
// `file`: 配置文件名
func Init(ctx context.Context, file string) error {
	logger := logrus.New().WithContext(ctx).WithField("func", "config.Init")

	logger.Infof("init config, config file: %v", file)
	// 从配置文件中获取相关配置
	if err := configor.New(&configor.Config{ENVPrefix: "OWL"}).Load(&appConfig, file); err != nil {
		logger.WithError(err).Error("cannot init config from config file: %v", file)
		return err
	}
	logger.Infof("inited config: %+v", appConfig)

	// 初始化 Web Socket Upgrader
	websocketUpgrader = websocket.Upgrader{
		HandshakeTimeout: time.Duration(appConfig.WebSocket.HandshakeTimeoutMills) * time.Millisecond,
		ReadBufferSize:   appConfig.WebSocket.ReadBuffer,
		WriteBufferSize:  appConfig.WebSocket.WriteBuffer,
		WriteBufferPool:  nil,
		Subprotocols:     appConfig.WebSocket.Subprotocols,
		Error: func(writer http.ResponseWriter, req *http.Request, status int, reason error) {
		},
		CheckOrigin:       func(r *http.Request) bool { return appConfig.WebSocket.CheckOrigin },
		EnableCompression: appConfig.WebSocket.EnableCompression,
	}
	logger.Infof("init web socket upgrader: %+v", websocketUpgrader)

	return nil
}
