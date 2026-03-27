package controllers

import (
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"gin-synolux/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	BaseController
	Service *service.ArticleService //依赖注入
}

// 获取首页文章(前3条)
func (c *ArticleController) HomeArticle(ctx *gin.Context) {
	//获取文䓬列表
	query := dto.ArticleQuery{
		Limit: 3,
		Status:   1,
	}
	list, _ := c.Service.List(query)

	c.SuccessJson(ctx, "success", gin.H{
		"list": list,
	})
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	catid, _ := strconv.Atoi(ctx.DefaultQuery("catid", "0"))
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
	query := dto.ArticleQuery{
		Title:    title,
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
func (c *ArticleController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	info, err := c.Service.GetById(id)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", info)
}

// 保存
func (c *ArticleController) Save(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	catid, _ := strconv.Atoi(ctx.PostForm("catid"))
	status, _ := strconv.Atoi(ctx.PostForm("status"))
	title := ctx.PostForm("title")
	info := ctx.PostForm("info")
	content := ctx.PostForm("content")
	author := ctx.PostForm("author")

	//组装实体
	entity := models.Article{
		Id:      id,
		Catid:   catid,
		Title:   title,
		Info:    info,
		Content: content,
		Author:  author,
		Status:  int8(status),
	}
	err := c.Service.Save(entity, false)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 删除
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.DeleteById(id, false)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 启用
func (c *ArticleController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.ChangeStatus(id, 1, false)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 禁用
func (c *ArticleController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.ChangeStatus(id, 0, false)
	if err != nil {
		c.handleError(ctx, err)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
