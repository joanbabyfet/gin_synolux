package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"strconv"

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
		IsAdmin:  true,
	}
	res, err := c.Service.List(query)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, res)
}

//获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	info, err := c.Service.GetById(id)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, info)
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
	data := models.Article{
		Id:      id,
		Catid:   catid,
		Title:   title,
		Info:    info,
		Content: content,
		Author:  author,
		Status:  int8(status),
		CreateUser: common.GetUserID(ctx),
		UpdateUser: common.GetUserID(ctx),
	}
	if err := c.Service.Save(data, true); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 删除
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	req := dto.ArticleDeleteReq{
		ID:     id,
		UserID: common.GetUserID(ctx),
		Role:   common.GetRole(ctx),
	}

	err := c.Service.DeleteById(req, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 启用
func (c *ArticleController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	req := dto.ArticleChangeStatusReq{
		ID:     id,
		Status: 1,
		UserID: common.GetUserID(ctx),
		Role:   common.GetRole(ctx),
	}

	err := c.Service.ChangeStatus(req, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 禁用
func (c *ArticleController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	req := dto.ArticleChangeStatusReq{
		ID:     id,
		Status: 0,
		UserID: common.GetUserID(ctx),
		Role:   common.GetRole(ctx),
	}

	err := c.Service.ChangeStatus(req, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}
