package main

import (
	"fmt"
	"log/slog"

	"IMProject/config"
	"IMProject/pkg/logger"
	"IMProject/pkg/mysql"
	"IMProject/routers"
	"IMProject/util"
)

func init() {
	//获取配置文件
	config.GetConfig("user_info")
	//初始化日志
	logger.InitLog(slog.Level(config.Cfg.LogLevel))
	//初始化rsa公钥和私钥
	util.InitRsaKey(config.Cfg.RsaAK)
	//初始化mysql连接池
	mysql.InitMysql(config.Cfg.Mysql)
}
func main() {
	routers.InitRouter()
	routers.InitUserInfo()
	routers.Routes.Run(fmt.Sprintf(":%d", config.Cfg.Port))
}
