package match

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 匹配相关的HTTP接口处理函数可以在这里定义
// 例如：
func MatchHandler(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Query("userId"), 10, 64)
	HandleMatchRequest(userId)
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "匹配请求已处理"})
}
