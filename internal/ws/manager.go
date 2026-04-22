package ws

import (
	"encoding/json"
	"game-server/internal/match"
	"game-server/pkg/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 统一游戏消息包结构（客户端↔服务器通用）
type GameMsg struct {
	Cmd  uint32 `json:"cmd"`  // 消息指令ID
	Data string `json:"data"` // 消息内容
	From uint64 `json:"from"` // 发送者ID
	To   uint64 `json:"to"`   // 接收者ID(私聊用)
}

// 消息指令定义
const (
	CmdWorldChat   = 1001 // 世界聊天
	CmdPrivateChat = 1002 // 私聊
)

// 单个玩家连接
type Client struct {
	Conn   *websocket.Conn
	UserId uint64
	Send   chan []byte // 消息发送管道
}

// 全局管理器：管理所有在线玩家
type Manager struct {
	clients sync.Map // 在线玩家 key:userId value:*Client
	mu      sync.RWMutex
}

var GlobalManager = &Manager{}

// 新增玩家连接
func (m *Manager) AddClient(client *Client) {
	m.clients.Store(client.UserId, client)
	logger.Info("玩家上线", zap.Uint64("userId", client.UserId))
}

// 移除玩家连接（掉线/退出）
func (m *Manager) RemoveClient(userId uint64) {
	m.clients.Delete(userId)
	logger.Info("玩家下线", zap.Uint64("userId", userId))
}

// 给指定玩家发消息
func (m *Manager) SendMsgToUser(userId uint64, msg []byte) {
	val, ok := m.clients.Load(userId)
	if !ok {
		return
	}
	client := val.(*Client)
	select {
	case client.Send <- msg:
	default:
		m.RemoveClient(userId)
	}
}

// 世界广播消息（所有人收到）
func (m *Manager) Broadcast(msg []byte) {
	m.clients.Range(func(key, value any) bool {
		client := value.(*Client)
		select {
		case client.Send <- msg:
		default:
			m.RemoveClient(key.(uint64))
		}
		return true
	})
}

// 读写消息循环
func (c *Client) ReadWriteLoop() {
	defer func() {
		GlobalManager.RemoveClient(c.UserId)
		c.Conn.Close()
	}()

	// 设置心跳超时 60秒无消息自动掉线
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// 读取客户端消息
	go func() {
		for {
			_, msgData, err := c.Conn.ReadMessage()
			if err != nil {
				return
			}

			logger.Info("收到玩家消息", zap.Uint64("userId", c.UserId), zap.ByteString("msg", msgData))

			// 解析统一消息协议
			var gameMsg GameMsg
			err = json.Unmarshal(msgData, &gameMsg)
			if err != nil {
				return
			}

			// 根据指令处理消息
			handleMessage(c, &gameMsg)
		}
	}()

	// 发送消息给客户端
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

// 处理消息
func handleMessage(client *Client, gameMsg *GameMsg) {
	switch gameMsg.Cmd {
	case CmdWorldChat:
		// 世界聊天
		handleWorldChat(client, gameMsg)
	case CmdPrivateChat:
		// 私聊
		handlePrivateChat(client, gameMsg)
	case 2000:
		match.HandleMatchRequest(client.UserId)
	}
}

// 处理世界聊天
func handleWorldChat(client *Client, gameMsg *GameMsg) {
	// 组装消息
	responseMsg := GameMsg{
		Cmd:  CmdWorldChat,
		From: client.UserId,
		Data: gameMsg.Data,
	}

	// 转json
	data, _ := json.Marshal(responseMsg)

	// 广播给所有在线玩家
	GlobalManager.Broadcast(data)

	logger.Info("世界聊天发送成功",
		zap.Uint64("sender", client.UserId),
		zap.String("content", gameMsg.Data),
	)
}

// 处理私聊
func handlePrivateChat(client *Client, gameMsg *GameMsg) {
	responseMsg := GameMsg{
		Cmd:  CmdPrivateChat,
		From: client.UserId,
		To:   gameMsg.To,
		Data: gameMsg.Data,
	}

	data, _ := json.Marshal(responseMsg)

	// 只发给目标玩家
	GlobalManager.SendMsgToUser(gameMsg.To, data)

	logger.Info("私聊发送成功",
		zap.Uint64("sender", client.UserId),
		zap.Uint64("target", gameMsg.To),
		zap.String("content", gameMsg.Data),
	)
}
