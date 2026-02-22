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
}

// 获取首页文章(前3条)
func (c *ArticleController) HomeArticle(ctx *gin.Context) {
	//获取文䓬列表
	query := dto.ArticleQuery{
		Limit: 3,
		Status:   1,
	}
	service_article := &service.ArticleService{}
	list, _ := service_article.List(query)
	
	//构造 JSON 的 map
	res := gin.H{
		"list": list,
	}

	c.SuccessJson(ctx, "success", res)
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	catid, _ := strconv.Atoi(ctx.Query("catid"))
	title := ctx.Query("title")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	page_size, _ := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || page_size < 1 {
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
	service_article := &service.ArticleService{}
	list, count := service_article.List(query)

	//构造 JSON 的 map
	res := gin.H{
		"list": list,
	}

	//显示总条数
	if query.Count {
		res["count"] = count
	}

	c.SuccessJson(ctx, "success", res)
}

// 获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	service_article := new(service.ArticleService)
	info, err := service_article.GetById(id)
	if err != nil {
		se := err.(*service.ServiceError);
		c.ErrorJson(ctx, se.Code, se.Msg, nil)
		return
	}

	c.SuccessJson(ctx, "success", info)
}

// 保存
func (c *ArticleController) Save(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	catid, _ := strconv.Atoi(ctx.PostForm("catid"))
	title := ctx.PostForm("title")
	info := ctx.PostForm("info")
	content := ctx.PostForm("content")
	author := ctx.PostForm("author")
	status, _ := strconv.Atoi(ctx.PostForm("status"))

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
	service_article := new(service.ArticleService)
	err := service_article.Save(entity, false)
	if err != nil {
		se := err.(*service.ServiceError);
		c.ErrorJson(ctx, se.Code, se.Msg, nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 删除
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	service_article := new(service.ArticleService)
	err := service_article.DeleteById(id, false)
	if err != nil {
		se := err.(*service.ServiceError);
		c.ErrorJson(ctx, se.Code, se.Msg, nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 启用
func (c *ArticleController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	service_article := new(service.ArticleService)
	err := service_article.ChangeStatus(id, 1, false)
	if err != nil {
		se := err.(*service.ServiceError);
		c.ErrorJson(ctx, se.Code, se.Msg, nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 禁用
func (c *ArticleController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	service_article := new(service.ArticleService)
	err := service_article.ChangeStatus(id, 0, false)
	if err != nil {
		se := err.(*service.ServiceError);
		c.ErrorJson(ctx, se.Code, se.Msg, nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
