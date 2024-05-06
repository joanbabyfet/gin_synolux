// 消费者, 队列需要单独控制台命令启动，与http服务独立避免相互影响
package main

import (
	"gin-synolux/jobs"
	"gin-synolux/queue"
	"gin-synolux/utils"

	"github.com/spf13/pflag"
)

var cfg = pflag.StringP("config", "", "", "apiserver config file path.")

func main() {
	//解析定义的标志
	pflag.Parse()

	//初始化配置
	if err := utils.InitConfig(*cfg); err != nil {
		panic(err)
	}

	subscriberMail := new(jobs.SubscribeMail) //实例化
	subscriberSMS := new(jobs.SubscribeSMS)
	forever := make(chan bool)

	q := queue.NewQueue()

	//队列执行的任务需要注册方可执行
	q.PushJob("mail", jobs.HandlerFunc(subscriberMail.ActionMail))
	q.PushJob("sms", jobs.HandlerFunc(subscriberSMS.ActionSMS))

	//提前规划好队列，可按延时时间来划分。可多个任务由一个队列来执行，也可以一个任务一个队列，一个队列可启动多个消费者
	go q.NewShareQueue("queue1")
	go q.NewShareQueue("queue2")

	defer q.Close()
	<-forever
}
