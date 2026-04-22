package antiCheat

import (
	"context"
	"fmt"
	"game-server/pkg/logger"
	"game-server/pkg/redis"
	"time"

	"go.uber.org/zap"
)

// 防作弊：接口频率限制 1秒最多5次请求
func CheckRateLimit(userId uint64) bool {
	key := fmt.Sprintf("anti:rate:%d", userId)
	count, err := redis.RDB.Incr(context.Background(), key).Result()
	if err != nil {
		return false
	}
	// 过期1秒
	if count == 1 {
		redis.RDB.Expire(context.Background(), key, time.Second)
	}
	// 超过5次拦截
	if count > 5 {
		logger.Warn("检测到玩家恶意刷接口作弊", zap.Uint64("userId", userId))
		return false
	}
	return true
}

// 防作弊：关键操作服务器强校验，拒绝客户端下发数值
func CheckScoreChange(oldScore, newScore int) bool {
	// 禁止单次加分超过100，防止客户端篡改金币积分
	if newScore-oldScore > 100 {
		logger.Warn("检测到积分异常篡改作弊")
		return false
	}
	return true
}
