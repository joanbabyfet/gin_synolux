package admin

import (
	"gin-synolux/common"
	"gin-synolux/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	AdminBaseController
	Service *service.UploadService
}

func NewUploadController(s *service.UploadService) *UploadController {
	return &UploadController{Service: s}
}

// 上传图片
func (c *UploadController) Upload(ctx *gin.Context) {
	thumb_w, _ := strconv.Atoi(ctx.PostForm("thumb_w")) //缩图宽
	thumb_h, _ := strconv.Atoi(ctx.PostForm("thumb_h")) //缩图高
	dir := ctx.DefaultPostForm("dir", "image")          //文件上传目录, 默认image
	file, err := ctx.FormFile("file")
	if err != nil {
		common.Fail(ctx, -1, "文件获取失败", nil)
		return
	}

	res, err := c.Service.Upload(file, dir, thumb_w, thumb_h)
	if err != nil {
		common.HandleError(ctx, err)
		return
	}
	common.Success(ctx, res)
}

// 下载文件
func (c *UploadController) Download(ctx *gin.Context) {
	filename := ctx.Query("filename") // 从请求拿, 比如 20240404/beaedec7c974a5c8e9a9f8770f9cec2b.png

	fullPath, _, err := c.Service.GetFilePath(filename)
	if err != nil {
		common.Fail(ctx, -1, err.Error(), nil)
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")              // 表示是文件流，唤起浏览器下载，一般设置了这个，就要设置文件名
	ctx.Header("Content-Disposition", "attachment; filename="+filename) // 用来指定下载下来的文件名
	ctx.Header("Content-Transfer-Encoding", "binary")                   // 表示传输过程中的编码形式，乱码问题可能就是因为它

	ctx.File(fullPath)
}
