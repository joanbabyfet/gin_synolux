package middleware

import (
	"gin-synolux/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		auth := ctx.GetHeader("Authorization")
		kv := strings.Split(auth, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			ctx.AbortWithStatusJSON(200, gin.H{
				"code": -1,
				"msg":  "未带token",
				"timestamp": common.Timestamp(),
			})
			return
		}

		payload, err := common.ValidateToken(kv[1])
		if err != nil {
			ctx.AbortWithStatusJSON(200, gin.H{
				"code": -2,
				"msg":  "未登录或登录超时",
				"timestamp": common.Timestamp(),
			})
			return
		}

		exists, _ := common.Redis.Exists("jwt:blacklist:" + kv[1]).Result()
		if exists > 0 {
			ctx.AbortWithStatusJSON(200, gin.H{
				"code": -3,
				"msg":  "token 已失效",
				"timestamp": common.Timestamp(),
			})
			return
		}

		// 存入 context
		ctx.Set("userID", payload.UserID)

		ctx.Next()
	}
}