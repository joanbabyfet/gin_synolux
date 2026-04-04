// 父控制器
package controllers

import (
	"gin-synolux/common"
	"gin-synolux/dto"
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
	var req dto.FeedbackSaveReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	req.CreateUser = "0"

	if err := c.Service.Save(&req, false); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}
