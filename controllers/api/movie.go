package controllers

import (
	"gin-synolux/dto"
	"gin-synolux/service"
	"gin-synolux/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieController struct {
	BaseController
	Service *service.MovieService //依赖注入
}

// 获取列表
func (c *MovieController) Index(ctx *gin.Context) {
	title := ctx.Query("title")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	if page < 1 {
		page = 1
	}

	if page_size < 1 {
		page_size = 10
	}

	//获取文䓬列表
	query := dto.MovieQuery{
		Title:    title,
		Pager: utils.Pager{
			Page:     page,
			PageSize: page_size,
		},
		Status:   1,
		Count: true,
	}
	list, count := c.Service.List(query)

	c.SuccessJson(ctx, "success", gin.H{
		"list":  list, 	//构造 JSON 的 map
		"count": count, //显示总条数
	})
}
