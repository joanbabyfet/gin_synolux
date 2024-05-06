package admin

import (
	"encoding/json"
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
	cache_key := "work_queue" //队列名
	data := map[string]interface{}{
		"to":   "+886958035350",
		"body": "测试",
	}
	for i := 1; i <= 2; i++ {
		job := utils.Job{Queue: "default", Task: "send_sms", Args: data}
		jobJSON, _ := json.Marshal(job)
		_, err := utils.Redis.RPush(cache_key, jobJSON).Result() //将一个或多个值插入到列表的尾部(最右边)
		if err != nil {
			panic(err)
		}
	}

	data = map[string]interface{}{
		"to":      "crwu0206@gmail.com",
		"subject": "测试",
		"body":    "测试测试测试",
	}
	for i := 1; i <= 10; i++ {
		job := utils.Job{Queue: "default", Task: "send_mail", Args: data}
		jobJSON, _ := json.Marshal(job)
		_, err := utils.Redis.RPush(cache_key, jobJSON).Result() //将一个或多个值插入到列表的尾部(最右边)
		if err != nil {
			panic(err)
		}
	}
	// 等待所有工作完成
	//time.Sleep(5 * time.Second)

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

	//创建一个新的RabbitMQ实例
	//forever := make(chan bool)
	// rabbitmq, err := utils.NewRabbitMQ("queue1", "", "", "amqp://guest:guest@localhost:5672/")
	// defer rabbitmq.Destroy()
	// if err != nil {
	// 	log.Println(err)
	// }

	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		rabbitmq.Publish("消息:" + strconv.Itoa(i))
	// 	}
	// }()

	// <-forever

	c.SuccessJson(ctx, "success", nil)
}
