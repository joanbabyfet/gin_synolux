package admin

import (
	"gin-synolux/common"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"strconv"

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
	query := dto.AdQuery{
		Title:    title,
		Catid:    catid,
		Page:     page,
		PageSize: page_size,
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
func (c *AdController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	info, err := c.Service.GetById(id)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, info)
}

// 保存
func (c *AdController) Save(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	catid, _ := strconv.Atoi(ctx.PostForm("catid"))
	status, _ := strconv.Atoi(ctx.PostForm("status"))
	sort, _ := strconv.Atoi(ctx.PostForm("sort"))
	title := ctx.PostForm("title")
	img := ctx.PostForm("img")
	url := ctx.PostForm("url")

	//组装实体
	data := models.Ad{
		Id:      id,
		Catid:   catid,
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
func (c *AdController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.DeleteById(id, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 启用
func (c *AdController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.ChangeStatus(id, 1, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}

// 禁用
func (c *AdController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	err := c.Service.ChangeStatus(id, 0, true)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, nil)
}