// 父控制器
package controllers

import (
	"gin-synolux/models"
	"gin-synolux/service"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
)

type FeedbackController struct {
	BaseController
}

// 保存
func (c *FeedbackController) Save(ctx *gin.Context) {
	name := ctx.PostForm("name")
	mobile := ctx.PostForm("mobile")
	email := ctx.PostForm("email")
	content := ctx.PostForm("content")

	//参数验证
	entity := models.Feedback{
		Name:    name,
		Mobile:  mobile,
		Email:   email,
		Content: content,
	}
	valid := validation.Validation{}
	valid.Required(entity.Name, "name")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	service_feedback := new(service.FeedbackService)
	stat, err := service_feedback.Save(entity)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
