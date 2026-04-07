package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/service"

	"github.com/gin-gonic/gin"
)

type AdController struct {
	AdminBaseController
	Service *service.AdService //依赖注入
}

func NewAdController(s *service.AdService) *AdController {
	return &AdController{Service: s}
}

// 获取列表
func (c *AdController) Index(ctx *gin.Context) {
	var req dto.AdQuery
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
func (c *AdController) Detail(ctx *gin.Context) {
	var req dto.AdDetailReq

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
func (c *AdController) Save(ctx *gin.Context) {
	var req dto.AdSaveReq

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
func (c *AdController) Delete(ctx *gin.Context) {
	var req dto.AdDeleteReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}
	//补充上下文
	req.UserID = common.GetUserID(ctx)
	req.Role = common.GetRole(ctx)

	err := c.Service.DeleteById(req, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 启用
func (c *AdController) Enable(ctx *gin.Context) {
	var req dto.AdChangeStatusReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}
	req.Status = 1
	req.UserID = common.GetUserID(ctx)
	req.Role = common.GetRole(ctx)

	err := c.Service.ChangeStatus(req, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 禁用
func (c *AdController) Disable(ctx *gin.Context) {
	var req dto.AdChangeStatusReq

	if err := ctx.ShouldBind(&req); err != nil {
		common.Fail(ctx, -1, common.GetValidMsg(err), nil)
		return
	}
	req.Status = 0
	req.UserID = common.GetUserID(ctx)
	req.Role = common.GetRole(ctx)

	err := c.Service.ChangeStatus(req, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}