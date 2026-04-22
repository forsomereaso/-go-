package api

import (
	"fmt"
	"game-server/internal/jwt"
	"game-server/pkg/mysql"
	"game-server/pkg/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 用户表结构体
type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Nickname string `gorm:"size:32"`
	Password string `gorm:"size:100"`
}

// 注册参数
type RegisterReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录参数
type LoginReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册接口
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//密码加密
	pwdHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := User{
		Nickname: req.Nickname,
		Password: string(pwdHash),
	}

	//写入数据库
	mysql.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "注册成功"})
}

// 登录接口
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//查询用户
	var user User
	mysql.DB.Where("nickname = ?", req.Nickname).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}

	//对比密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "密码错误"})
		return
	}

	//生成token
	token, _ := jwt.GenerateToken(user.ID)

	//redis保存在线状态
	redis.Set(fmt.Sprintf("online:%d", user.ID), token, time.Hour*24)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "登录成功",
		"token": token,
	})
}
