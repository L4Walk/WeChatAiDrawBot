package Log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger

func init() {
	// 配置日志
	// 使用 lumberjack 实现日志自动分片
	lumberjackLogger := &lumberjack.Logger{
		Filename: getLogFilename(), // 获取今天的日志文件名
		//MaxSize:    100,              // 每个日志文件的最大尺寸（单位：MB）
		//MaxBackups: 180,              // 最多保留的旧日志文件个数
		//MaxAge:     180,              // 保留旧日志文件的最大天数（设置为1，即保留当天的日志文件）
		Compress: true, // 是否压缩旧日志文件
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	// 创建Core，使用 lumberjack 实现的 WriteSyncer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		zap.InfoLevel,
	)

	var err error
	Logger = zap.New(core)

	if err != nil {
		panic("无法初始化Zap Logger：" + err.Error())
	}

	// 确保在程序启动时将日志缓冲区中的内容写入文件
	defer Logger.Sync()

	// 每天零点更新日志文件名
	go updateLogFile(lumberjackLogger)
}

// GetLogger 获取全局的Zap Logger实例
func GetLogger() *zap.Logger {
	return Logger
}

// 每天零点更新日志文件名
func updateLogFile(logger *lumberjack.Logger) {
	for {
		now := time.Now()
		next := now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		duration := next.Sub(now)
		<-time.After(duration)
		logger.Filename = "logs/app_" + now.Format("2006-01-02") + ".log"
	}
}

// 获取今天的日志文件名，确保启动时也能记录
func getLogFilename() string {
	now := time.Now()
	filename := fmt.Sprintf("logs/app_%s.log", now.Format("2006-01-02"))
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 如果今天的日志文件不存在，则创建
		_, _ = os.Create(filename)
	}
	return filename
}
