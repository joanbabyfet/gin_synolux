package controllers

import (
	"gin-synolux/common"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	BaseController
}

// 测试用
func (c *TestController) Test(ctx *gin.Context) {
	common.Success(ctx, nil)
}
