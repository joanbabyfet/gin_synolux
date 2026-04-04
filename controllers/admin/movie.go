package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"

	"github.com/gin-gonic/gin"
)

type MovieController struct {
	AdminBaseController
	Service *service.MovieService //依赖注入
}

func NewMovieController(s *service.MovieService) *MovieController {
	return &MovieController{Service: s}
}

// 获取列表
func (c *MovieController) Index(ctx *gin.Context) {
	var req dto.MovieQuery

	if err := ctx.ShouldBind(&req); err != nil {
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
func (c *MovieController) Detail(ctx *gin.Context) {
	var req dto.MovieDetailReq

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

// 保存
func (c *MovieController) Save(ctx *gin.Context) {
	var req dto.MovieSaveReq

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
func (c *MovieController) Delete(ctx *gin.Context) {
	var req dto.MovieDeleteReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
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
func (c *MovieController) Enable(ctx *gin.Context) {
	var req dto.MovieChangeStatusReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
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
func (c *MovieController) Disable(ctx *gin.Context) {
	var req dto.MovieChangeStatusReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, "参数错误", nil)
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