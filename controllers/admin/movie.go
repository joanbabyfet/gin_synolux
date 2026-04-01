package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"strconv"

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
	query := dto.MovieQuery{
		Title:    title,
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
func (c *MovieController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	info, err := c.Service.GetById(id)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, info)
}

// 保存
func (c *MovieController) Save(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	status, _ := strconv.Atoi(ctx.PostForm("status"))
	sort, _ := strconv.Atoi(ctx.PostForm("sort"))
	title := ctx.PostForm("title")
	img := ctx.PostForm("img")
	url := ctx.PostForm("url")

	//组装实体
	data := models.Movie{
		Id:      id,
		Title:   title,
		Img:     img,
		Url: 	 url,
		Sort: 	 int16(sort),
		Status:  int8(status),
	}
	if err := c.Service.Save(data, true); err != nil {
		common.HandleError(ctx, err)
		return
	}

	common.Success(ctx, nil)
}

// 删除
func (c *MovieController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	req := dto.MovieDeleteReq{
		ID:     id,
		UserID: common.GetUserID(ctx),
		Role:   common.GetRole(ctx),
	}

	err := c.Service.DeleteById(req, false)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 启用
func (c *MovieController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	req := dto.MovieChangeStatusReq{
		ID:     id,
		Status: 1,
		UserID: common.GetUserID(ctx),
		Role:   common.GetRole(ctx),
	}

	err := c.Service.ChangeStatus(req, false)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 禁用
func (c *MovieController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	req := dto.MovieChangeStatusReq{
		ID:     id,
		Status: 1,
		UserID: common.GetUserID(ctx),
		Role:   common.GetRole(ctx),
	}

	err := c.Service.ChangeStatus(req, false)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}