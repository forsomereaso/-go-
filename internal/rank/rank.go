package rank

import (
	"context"
	"game-server/pkg/redis"
	"strconv"

	redisv8 "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// 排行榜Key
const rankKey = "game:rank:score"

// AddScore 给玩家增加积分
func AddScore(userID uint64, score int) {
	redis.RDB.ZAdd(ctx, rankKey, &redisv8.Z{
		Score:  float64(score),
		Member: userID,
	})
}

// GetTop100 获取前100名
func GetTop100() ([]redisv8.Z, error) {
	return redis.RDB.ZRevRangeWithScores(ctx, rankKey, 0, 99).Result()
}

// GetUserRank 获取玩家排名 & 分数
func GetUserRank(userID uint64) (rank int64, score float64, err error) {
	// 获取排名（从0开始，所以+1）
	rank, err = redis.RDB.ZRevRank(ctx, rankKey, strconv.FormatUint(userID, 10)).Result()
	if err != nil {
		return 0, 0, err
	}
	rank += 1

	// 获取分数
	score, err = redis.RDB.ZScore(ctx, rankKey, strconv.FormatUint(userID, 10)).Result()
	return rank, score, err
}
