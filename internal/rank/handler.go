package rank

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddScoreHandler 增加积分接口
func AddScoreHandler(c *gin.Context) {
	uidStr := c.Query("uid")
	scoreStr := c.Query("score")

	uid, _ := strconv.ParseUint(uidStr, 10, 64)
	score, _ := strconv.Atoi(scoreStr)

	AddScore(uid, score)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "积分增加成功",
	})
}

// GetRankHandler 获取排行榜
func GetRankHandler(c *gin.Context) {
	list, err := GetTop100()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": list,
	})
}

// GetUserRankHandler 获取自己排名
func GetUserRankHandler(c *gin.Context) {
	uidStr := c.Query("uid")
	uid, _ := strconv.ParseUint(uidStr, 10, 64)

	rank, score, err := GetUserRank(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "未上榜"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"rank":  rank,
		"score": score,
	})
}
