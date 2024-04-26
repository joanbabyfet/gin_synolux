package controllers

import (
	"gin-synolux/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UploadController struct {
	BaseController
}

// 上传图片
func (c *UploadController) Upload(ctx *gin.Context) {
	f, err := ctx.FormFile("filename")                  //获取文件信息
	dir := ctx.PostForm("dir")                          //文件上传目录, 默认image
	thumb_w, _ := strconv.Atoi(ctx.PostForm("thumb_w")) //缩图宽
	thumb_h, _ := strconv.Atoi(ctx.PostForm("thumb_h")) //缩图高
	if err != nil {
		c.ErrorJson(ctx, -1, "上传文件失败", nil)
		return
	}
	if dir == "" {
		dir = "image"
	}

	service_upload := new(service.UploadService)
	stat, data, err := service_upload.Upload(ctx, f, dir, thumb_w, thumb_h)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", data)
}

// 下载文件
func (c *UploadController) Download(ctx *gin.Context) {
	//获取文件名 20240404/beaedec7c974a5c8e9a9f8770f9cec2b.png
	filename := "20240423/5bbbfa7e754caf3c16b4cc4a774f35b8.jpg"
	ctx.Header("Content-Type", "application/octet-stream")              // 表示是文件流，唤起浏览器下载，一般设置了这个，就要设置文件名
	ctx.Header("Content-Disposition", "attachment; filename="+filename) // 用来指定下载下来的文件名
	ctx.Header("Content-Transfer-Encoding", "binary")                   // 表示传输过程中的编码形式，乱码问题可能就是因为它
	upload_dir := viper.GetString("upload_dir")
	ctx.File(upload_dir + "/image/" + filename)
}
