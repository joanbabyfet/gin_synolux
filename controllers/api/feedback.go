// 父控制器
package controllers

import (
	"gin-synolux/common"
	"gin-synolux/models"
	"gin-synolux/service"

	"github.com/gin-gonic/gin"
)

type FeedbackController struct {
	BaseController
	Service *service.FeedbackService //依赖注入
}

func NewFeedbackController(s *service.FeedbackService) *FeedbackController {
	return &FeedbackController{Service: s}
}

// 保存
func (c *FeedbackController) Save(ctx *gin.Context) {
	name := ctx.PostForm("name")
	mobile := ctx.PostForm("mobile")
	email := ctx.PostForm("email")
	content := ctx.PostForm("content")

	//参数验证
	data := models.Feedback{
		Name:    name,
		Mobile:  mobile,
		Email:   email,
		Content: content,
	}
	err := c.Service.Save(data, false)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}
