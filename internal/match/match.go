package match

import (
	"context"
	"encoding/json"
	"fmt"
	"game-server/internal/msg"
	"game-server/pkg/kafka"
	"game-server/pkg/logger"
	"game-server/pkg/redis"

	redisv8 "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var ctx = context.Background()

const matchKey = "game:match:queue" // 分布式匹配队列

// 玩家加入匹配队列
func JoinMatch(userId uint64, score int) {
	// 以玩家积分为分数排序，同分段优先匹配
	redis.RDB.ZAdd(ctx, matchKey, &redisv8.Z{
		Score:  float64(score),
		Member: userId,
	})
	logger.Info("玩家进入匹配队列", zap.Uint64("userId", userId), zap.Int("score", score))
}

// 定时匹配检测（每1秒扫描队列，找到2人就匹配成功）
func TickMatch() {
	// 取出匹配队列前2个玩家
	list, err := redis.RDB.ZRangeWithScores(ctx, matchKey, 0, 1).Result()
	if err != nil || len(list) < 2 {
		return
	}

	// 拿到两个玩家ID
	p1, _ := list[0].Member.(uint64)
	p2, _ := list[1].Member.(uint64)

	// 从匹配队列移除（匹配成功）
	redis.RDB.ZRem(ctx, matchKey, p1, p2)

	// 发送匹配成功消息到Kafka
	kafka.Send("match_success", fmt.Sprintf("p1:%d p2:%d", p1, p2))
	logger.Info("匹配成功", zap.Uint64("p1", p1), zap.Uint64("p2", p2))

	// 组装匹配成功消息
	roomMsg := msg.GameMsg{
		Cmd:  2001, // 匹配成功指令
		Data: "匹配成功，进入房间",
		From: p1,
		To:   p2,
	}
	data, _ := json.Marshal(roomMsg)
	// 发送Kafka消息通知匹配成功
	kafka.Send("match_success", string(data))

	// 通过Kafka消息通知，而不是直接调用WebSocket
	// 这样可以避免循环导入问题
	// 实际项目中，WebSocket服务可以订阅Kafka的match_success主题来处理匹配成功的通知
}

// 处理玩家发起匹配请求
func HandleMatchRequest(userId uint64) {
	// 获取玩家当前积分作为匹配分
	scoreStr, _ := redis.Get("user:score:" + fmt.Sprintf("%d", userId))
	score := 1000 // 默认分数
	if scoreStr != "" {
		_, err := fmt.Sscanf(scoreStr, "%d", &score)
		if err != nil {
			score = 1000
		}
	}
	JoinMatch(userId, score)
}
