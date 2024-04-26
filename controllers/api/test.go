package controllers

import (
	"github.com/gin-gonic/gin"
)

type TestController struct {
	BaseController
}

// 测试用
func (c *TestController) Test(ctx *gin.Context) {
	c.SuccessJson(ctx, "success", nil)
}
