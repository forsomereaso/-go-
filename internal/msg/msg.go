package msg

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
