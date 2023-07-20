package Log

import (
	"io/ioutil"
	"log"
	"os"
)

var Logger *log.Logger
var LogFile *os.File

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func LogInit() {
	// 创建或打开用于日志输出的文件
	var err error
	LogFile, err = os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	Trace = log.New(ioutil.Discard,
		"Trace: ", log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

	//defer LogFile.Close()

	// 设置日志输出位置为创建的文件——启用不会在控制台中输出
	log.SetOutput(LogFile)

	Trace.Println("Hello L4Walk!")
	// 记录日志
	Info.Println("This is a log message.")

	defer LogFile.Sync()
}

/*
func LogInit2() {
	var err error
	LogFile, err = os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// 创建一个WriteSyncer，将日志输出到文件
	fileWriteSyncer := zapcore.AddSync(LogFile)

	// 配置日志编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式使用ISO8601

	// 创建Core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器
		fileWriteSyncer,                          // 输出位置
		zap.InfoLevel,                            // 日志级别
	)

	// 创建Logger
	Logger = zap.New(core)

	// 确保缓存区中的日志都被写入
	defer Logger.Sync()

	// 全局调用logger
	zap.ReplaceGlobals(Logger)

	defer LogFile.Close()
}

func GetLogger() *zap.Logger {
	return Logger
}
*/
