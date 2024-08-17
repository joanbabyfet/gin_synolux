package controllers

import (
	"gin-synolux/dto"
	"gin-synolux/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieController struct {
	BaseController
}

// 获取列表
func (c *MovieController) Index(ctx *gin.Context) {
	title := ctx.Query("title")
	page, _ := strconv.Atoi(ctx.Query("page"))
	page_size, _ := strconv.Atoi(ctx.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if page_size < 1 {
		page_size = 10
	}

	//获取文䓬列表
	query := dto.MovieQuery{}
	query.Title = title
	query.Page = page
	query.PageSize = page_size
	query.Status = 1
	service_movie := new(service.MovieService)
	list, total := service_movie.PageList(query)

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["total"] = total
	resp["list"] = list
	c.SuccessJson(ctx, "success", resp)
}
