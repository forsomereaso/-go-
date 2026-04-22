package ws

import (
	"game-server/internal/jwt"
	"game-server/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"code": 1, "msg": "请先登录，缺少token"})
		return
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"code": 1, "msg": "token无效，请重新登录"})
		return
	}
	userId := claims.UserId

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("websocket升级失败", zap.Error(err))
		return
	}

	client := &Client{
		Conn:   conn,
		UserId: userId,
		Send:   make(chan []byte, 256),
	}

	GlobalManager.AddClient(client)
	client.ReadWriteLoop()
}
