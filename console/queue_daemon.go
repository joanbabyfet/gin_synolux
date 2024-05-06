// 消费者, 队列需要单独控制台命令启动，与http服务独立避免相互影响
package main

import (
	"gin-synolux/utils"
	"log"
)

func main() {
	forever := make(chan bool) //注册在主进程，不需要阻塞

	//创建一个新的RabbitMQ实例
	rabbitmq, err := utils.NewRabbitMQ("queue1", "", "", "amqp://guest:guest@localhost:5672/")
	defer rabbitmq.Destroy()
	if err != nil {
		log.Println(err)
	}

	// 执行消费
	go func() {
		msgs, err3 := rabbitmq.Consume()
		if err3 != nil {
			log.Println(err3)
		}
		for d := range msgs {
			log.Printf("接受到了：%s", string(d.Body))
		}
	}()

	<-forever
}
