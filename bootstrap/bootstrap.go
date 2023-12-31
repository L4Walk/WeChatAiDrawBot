package bootstrap

import (
	"WeChatAiDrawBot/Database"
	"WeChatAiDrawBot/Log"
	"WeChatAiDrawBot/config"
	"WeChatAiDrawBot/handlers"
	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

func LogInit() {
	Logger = Log.GetLogger()
}

func DatabaseInit() {
	// 连接数据库
	db, err := Database.OpenDBConnetction(config.LoadConfig().MYSQL_CONNECT_STRING)

	if err != nil {
		Logger.Fatal("Error connecting to the database: ", zap.Error(err))
	}

	// 确保数据库能断开连接
	defer db.Close()

	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		Logger.Fatal("Error connecting to the database: ", zap.Error(err))
	}

	Log.Logger.Info("Connected to the database successfully!")
}

func Run() {
	// 连接日志
	LogInit()

	// 连接数据库
	DatabaseInit()

	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
