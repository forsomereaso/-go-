package chat

import (
	"encoding/json"
	"game-server/internal/msg"
	"game-server/internal/ws"
	"game-server/pkg/logger"

	"go.uber.org/zap"
)

// 处理世界聊天消息
func HandleWorldChat(fromUserId uint64, content string) {
	// 组装消息
	gameMsg := msg.GameMsg{
		Cmd:  msg.CmdWorldChat,
		From: fromUserId,
		Data: content,
	}

	// 转json
	data, _ := json.Marshal(gameMsg)

	// 广播给所有在线玩家
	ws.GlobalManager.Broadcast(data)

	logger.Info("世界聊天发送成功",
		zap.Uint64("sender", fromUserId),
		zap.String("content", content),
	)
}

// 处理私聊消息
func HandlePrivateChat(fromUserId, toUserId uint64, content string) {
	gameMsg := msg.GameMsg{
		Cmd:  msg.CmdPrivateChat,
		From: fromUserId,
		To:   toUserId,
		Data: content,
	}

	data, _ := json.Marshal(gameMsg)

	// 只发给目标玩家
	ws.GlobalManager.SendMsgToUser(toUserId, data)

	logger.Info("私聊发送成功",
		zap.Uint64("sender", fromUserId),
		zap.Uint64("target", toUserId),
		zap.String("content", content),
	)
}
