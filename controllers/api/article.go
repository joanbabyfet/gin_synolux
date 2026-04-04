package controllers

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	BaseController
	Service *service.ArticleService //依赖注入
}

func NewArticleController(s *service.ArticleService) *ArticleController {
	return &ArticleController{Service: s}
}

// 获取首页文章(前3条)
func (c *ArticleController) HomeArticle(ctx *gin.Context) {
	status := int8(1)
	//获取文䓬列表
	query := dto.ArticleQuery{
		Limit: 3,
		Status:   &status,
	}
	res, err := c.Service.List(query)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, res)
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	catid, _ := strconv.Atoi(ctx.DefaultQuery("catid", "0"))
	title := ctx.Query("title")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	status := int8(1)

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
		Page:     page,
		PageSize: page_size,
		Status:   &status,
		Count:    true,
	}
	res, err := c.Service.List(query)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, res)
}

// 获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	var req dto.ArticleDetailReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}

	info, err := c.Service.GetById(req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, info)
}
