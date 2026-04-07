package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	AdminBaseController
	Service *service.ArticleService //依赖注入
}

func NewArticleController(s *service.ArticleService) *ArticleController {
	return &ArticleController{Service: s}
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	var req dto.ArticleQuery

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
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

//获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	var req dto.ArticleDetailReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}

	info, err := c.Service.GetById(req)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, info)
}

// 保存
func (c *ArticleController) Save(ctx *gin.Context) {
	var req dto.ArticleSaveReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}

	req.CreateUser = common.GetUserID(ctx)
	req.UpdateUser = common.GetUserID(ctx)

	if err := c.Service.Save(&req, true); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 删除
func (c *ArticleController) Delete(ctx *gin.Context) {
	var req dto.ArticleDeleteReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}

	req.UserID = common.GetUserID(ctx)
	req.Role = common.GetRole(ctx)

	if err := c.Service.DeleteById(req, true); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 启用
func (c *ArticleController) Enable(ctx *gin.Context) {
	var req dto.ArticleChangeStatusReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}
	req.Status = 1
	req.UserID = common.GetUserID(ctx)
	req.Role = common.GetRole(ctx)

	if err := c.Service.ChangeStatus(req, true); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 禁用
func (c *ArticleController) Disable(ctx *gin.Context) {
	var req dto.ArticleChangeStatusReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}
	req.Status = 0
	req.UserID = common.GetUserID(ctx)
	req.Role = common.GetRole(ctx)

	if err := c.Service.ChangeStatus(req, true); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}
