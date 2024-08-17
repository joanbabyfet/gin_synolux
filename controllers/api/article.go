package controllers

import (
	"encoding/json"
	"gin-synolux/dto"
	"gin-synolux/models"
	"gin-synolux/service"
	"gin-synolux/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/thedevsaddam/govalidator"
)

type ArticleController struct {
	BaseController
}

// 获取首页文章
func (c *ArticleController) HomeArticle(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	//获取文章列表
	query := dto.ArticleQuery{}
	query.Limit = limit
	query.Status = 1
	service_article := new(service.ArticleService)
	list := service_article.All(query)

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["list"] = list
	c.SuccessJson(ctx, "success", resp)
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	catid, _ := strconv.Atoi(ctx.Query("catid"))
	title := ctx.Query("title")
	page, _ := strconv.Atoi(ctx.Query("page"))
	page_size, _ := strconv.Atoi(ctx.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if page_size < 1 {
		page_size = 10
	}

	//获取文䓬列表
	query := dto.ArticleQuery{}
	query.Title = title
	query.Catid = catid
	query.Page = page
	query.PageSize = page_size
	query.Status = 1
	service_article := new(service.ArticleService)
	list, total := service_article.PageList(query)

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["total"] = total
	resp["list"] = list
	c.SuccessJson(ctx, "success", resp)
}

// 获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	//参数验证
	entity := models.Article{Id: id}
	rules := govalidator.MapData{}
	rules["id"] = []string{"required"}
	messages := govalidator.MapData{}
	messages["id"] = []string{"required:id 不能为空"}
	opts := govalidator.Options{
		Data:            &entity,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	valid := govalidator.New(opts)
	e := valid.ValidateStruct()
	if len(e) > 0 {
		for _, err := range e {
			c.ErrorJson(ctx, -1, err[0], nil)
			return
		}
	}

	//redis缓存
	info := new(models.Article)
	cache_key := "article:id:" + strconv.Itoa(id)
	v, err := utils.Redis.Get(cache_key).Result()
	if err == redis.Nil {
		//redis不存在则跟库拿
		service_article := new(service.ArticleService)
		info, err = service_article.GetById(id)
		if err != nil {
			c.ErrorJson(ctx, -2, err.Error(), nil)
			return
		}
		str, _ := json.Marshal(&info)                //struct转成json字符串, 返回[]byte
		utils.Redis.SetNX(cache_key, string(str), 0) //永不过期
	} else {
		json.Unmarshal([]byte(v), &info) //json字符串转成struct
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

	//参数验证
	entity := models.Article{
		Id:      id,
		Catid:   catid,
		Title:   title,
		Info:    info,
		Content: content,
		Author:  author,
		Status:  int8(status),
	}
	rules := govalidator.MapData{}
	messages := govalidator.MapData{}
	if entity.Id > 0 {
		rules["id"] = []string{"required"}
		messages["id"] = []string{"required:id 不能为空"}
	}
	rules["title"] = []string{"required"}
	messages["title"] = []string{"required:title 不能为空"}
	opts := govalidator.Options{
		Data:            &entity,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	valid := govalidator.New(opts)
	e := valid.ValidateStruct()
	if len(e) > 0 {
		for _, err := range e {
			c.ErrorJson(ctx, -1, err[0], nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.Save(entity)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 删除
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//参数验证
	entity := models.Article{Id: id}
	rules := govalidator.MapData{}
	messages := govalidator.MapData{}
	rules["id"] = []string{"required"}
	messages["id"] = []string{"required:id 不能为空"}
	opts := govalidator.Options{
		Data:            &entity,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	valid := govalidator.New(opts)
	e := valid.ValidateStruct()
	if len(e) > 0 {
		for _, err := range e {
			c.ErrorJson(ctx, -1, err[0], nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.DeleteById(id)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 启用
func (c *ArticleController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//参数验证
	entity := models.Article{Id: id}
	rules := govalidator.MapData{}
	messages := govalidator.MapData{}
	rules["id"] = []string{"required"}
	messages["id"] = []string{"required:id 不能为空"}
	opts := govalidator.Options{
		Data:            &entity,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	valid := govalidator.New(opts)
	e := valid.ValidateStruct()
	if len(e) > 0 {
		for _, err := range e {
			c.ErrorJson(ctx, -1, err[0], nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.EnableById(id)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}

// 禁用
func (c *ArticleController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//参数验证
	entity := models.Article{Id: id}
	rules := govalidator.MapData{}
	messages := govalidator.MapData{}
	rules["id"] = []string{"required"}
	messages["id"] = []string{"required:id 不能为空"}
	opts := govalidator.Options{
		Data:            &entity,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	valid := govalidator.New(opts)
	e := valid.ValidateStruct()
	if len(e) > 0 {
		for _, err := range e {
			c.ErrorJson(ctx, -1, err[0], nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.DisableById(id)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
