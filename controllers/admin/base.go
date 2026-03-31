// 父控制器
package admin

import (
	"github.com/gin-gonic/gin"
)

type AdminBaseController struct {
}

// 初始化, 先于Prepare函数
func init() {
}

// 定义prepare方法, 用户扩展用
func (c *AdminBaseController) Prepare() {
}

// 获取客户端ip
func (c *AdminBaseController) getClientIp(ctx *gin.Context) string {
	return ctx.ClientIP()
}