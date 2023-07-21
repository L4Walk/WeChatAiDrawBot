package main

import (
	"WeChatAiDrawBot/Log"
	"WeChatAiDrawBot/bootstrap"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func main() {
	Logger = Log.GetLogger()
	bootstrap.Run()
}
