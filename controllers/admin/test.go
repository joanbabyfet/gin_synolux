package admin

import (
	"encoding/json"
	"gin-synolux/jobs"
	"gin-synolux/queue"
	"gin-synolux/utils"

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

	// 生产者, 将工作添加到工作队列
	data := map[string]interface{}{
		"to":   "+886912345678",
		"body": "测试",
	}
	for i := 1; i <= 2; i++ {
		job := utils.Job{Queue: "queue1", Action: "sms", Args: data}
		jobJSON, _ := json.Marshal(job)
		_, err := utils.Redis.RPush(job.Queue, jobJSON).Result() //将一个或多个值插入到列表的尾部(最右边)
		if err != nil {
			panic(err)
		}
	}

	data = map[string]interface{}{
		"to":      "example@example.com",
		"subject": "测试",
		"body":    "测试测试测试",
	}
	for i := 1; i <= 10; i++ {
		job := utils.Job{Queue: "queue1", Action: "mail", Args: data}
		jobJSON, _ := json.Marshal(job)
		_, err := utils.Redis.RPush(job.Queue, jobJSON).Result() //将一个或多个值插入到列表的尾部(最右边)
		if err != nil {
			panic(err)
		}
	}

	//base64加密
	// str := "123456"
	// encode_str := base64.StdEncoding.EncodeToString([]byte(str))
	// fmt.Println("加密结果：", encode_str)

	//base64解密
	// decode_byte, _ := base64.StdEncoding.DecodeString(encode_str)
	// fmt.Println("解密结果：", string(decode_byte))

	// Hash 加密 MD5
	// hash := md5.Sum([]byte(str))
	// encode_str = fmt.Sprintf("%x", hash)
	// fmt.Println("加密结果：", encode_str)

	// 对称加密
	// key := []byte("1234567890123456") // 密钥必须是 16 位
	// ciphertext, _ := utils.Encrypt(key, []byte(str))
	// fmt.Println("加密结果：", string(ciphertext))

	// 对称解密
	// plaintext, _ := utils.Decrypt(key, ciphertext)
	// fmt.Println("解密结果：", string(plaintext))

	c.SuccessJson(ctx, "success", nil)
}

// 测试生产者
func (c *TestController) Queue(ctx *gin.Context) {
	queue.NewSender("queue1", "mail", jobs.SubscribeMail{
		To: "example@example.com", Subject: "测试", Body: "测试测试测试"}).Send()
	queue.NewSender("queue2", "sms", jobs.SubscribeSMS{
		To: "+886912345678", Body: "短信测试"}).Send()

	c.SuccessJson(ctx, "success", nil)
}
