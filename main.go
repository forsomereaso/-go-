package main

import (
	"game-server/api"
	"game-server/internal/rank"
	"game-server/internal/ws"
	"game-server/pkg/config"
	"game-server/pkg/kafka"
	"game-server/pkg/logger"
	"game-server/pkg/mysql"
	"game-server/pkg/redis"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	//1.加载配置
	if err := config.Init(); err != nil {
		panic(err)
	}

	//2.初始化日志
	logger.Init(&config.Conf.Log)

	//3.初始化MySQL
	mysql.Init()

	//4.初始化Redis
	redis.Init()

	//5.初始化Kafka
	kafka.Init()

	//自动迁移用户表
	mysql.DB.AutoMigrate(&api.User{})

	//6.启动Gin网关
	r := gin.Default()

	// HTTP登录注册接口
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", api.Register)
		v1.POST("/login", api.Login)
	}

	// 排行榜接口
	rankGroup := r.Group("/api/rank")
	{
		rankGroup.GET("/add", rank.AddScoreHandler)
		rankGroup.GET("/list", rank.GetRankHandler)
		rankGroup.GET("/user", rank.GetUserRankHandler)
	}

	// WebSocket游戏长连接接口
	r.GET("/ws", ws.WsHandler)

	logger.Info("游戏服务器全部启动成功", zap.Int("端口", config.Conf.App.Port))

	//启动服务
	r.Run(":8080")
}
