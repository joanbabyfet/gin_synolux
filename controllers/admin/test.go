package admin

import (
	"github.com/gin-gonic/gin"
)

type TestController struct {
	AdminBaseController
}

// 测试用
func (c *TestController) Test(ctx *gin.Context) {
	//发送邮件
	// ok := utils.SendMail("example@gmail.com", "测试", "测试测试测试")
	// if !ok {
	// 	c.ErrorJson(ctx, -1, "发送失败", nil)
	// 	return
	// }

	//发送短信
	// ok := utils.SendSMS("+886912345678", "测试")
	// if !ok {
	// 	c.ErrorJson(ctx, -1, "发送失败", nil)
	// 	return
	// }

	//fmt.Println(utils.DateToUnix("2024-04-25 08:47:00"))

	//多语言
	//fmt.Println(ginI18n.MustGetMessage(ctx, "api_param_error"))

	//加载html
	//ctx.HTML(http.StatusOK, "index.html", nil)
	//fmt.Println(viper.GetString("runmode"))
	//fmt.Println(utils.UniqueId())

	c.SuccessJson(ctx, "success", nil)
}
