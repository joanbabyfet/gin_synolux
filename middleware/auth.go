package middleware

import (
	"gin-synolux/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		auth := ctx.GetHeader("Authorization")
		kv := strings.Split(auth, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			common.Fail(ctx, -1, "未带token", nil)
			ctx.Abort() //阻止后续 handler 执行
			return
		}

		payload, err := common.ValidateToken(kv[1])
		if err != nil {
			common.Fail(ctx, -2, "未登录或登录超时", nil)
			ctx.Abort()
			return
		}

		// 权限判断
        if payload.Role != requiredRole {
			common.Fail(ctx, -3, "无权限", nil)
			ctx.Abort()
            return
        }

		exists, _ := common.Redis.Exists("jwt:blacklist:" + kv[1]).Result()
		if exists > 0 {
			common.Fail(ctx, -4, "token已失效", nil)
			ctx.Abort()
			return
		}

		// 存入 context
		ctx.Set("userID", payload.UserID)
		ctx.Set("role", payload.Role)

		ctx.Next()
	}
}