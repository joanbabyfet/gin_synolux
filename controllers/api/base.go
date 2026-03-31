// 父控制器
package controllers

import (
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

// 初始化, 先于Prepare函数
func init() {
}

// 定义prepare方法, 用户扩展用
func (c *BaseController) Prepare() {
}

// 获取客户端ip
func (c *BaseController) getClientIp(ctx *gin.Context) string {
	return ctx.ClientIP()
}