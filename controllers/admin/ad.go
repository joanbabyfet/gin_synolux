package admin

import (
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"gin-synolux/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdController struct {
	AdminBaseController
	Service *service.AdService
}

// 获取列表
func (c *AdController) Index(ctx *gin.Context) {
	catid, _ := strconv.Atoi(ctx.DefaultQuery("catid", "0"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	if page < 1 {
		page = 1
	}
	if page_size < 1 {
		page_size = 10
	}
	
	//获取广告列表
	query := dto.AdQuery{
		Catid:    catid,
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

// 获取详情
func (c *AdController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))
	
	info, err := c.Service.GetById(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", info)
}

// 保存
func (c *AdController) Save(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	catid, _ := strconv.Atoi(ctx.PostForm("catid"))
	title := ctx.PostForm("title")
	status, _ := strconv.Atoi(ctx.PostForm("status"))

	//组装实体
	entity := models.Ad{
		Id:     id,
		Catid:  catid,
		Title:  title,
		Status: int8(status),
	}
	err := c.Service.Save(entity, true)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 删除
func (c *AdController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.DeleteById(id, true)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 启用
func (c *AdController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.ChangeStatus(id, 1, true)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 禁用
func (c *AdController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.ChangeStatus(id, 0, true)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}