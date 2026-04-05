package controllers

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"

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
	var req dto.ArticleQuery

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}
	//补充上下文
	req.Limit = 3
	req.Count = false
	v := int8(1)
	req.Status = &v

	res, err := c.Service.List(req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, res)
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	var req dto.ArticleQuery

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
		return
	}
	//补充上下文
	req.Count = true

	res, err := c.Service.List(req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, res)
}

// 获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	var req dto.ArticleDetailReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
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
