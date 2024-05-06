package main

// import (
// 	"gin-synolux/utils"
// 	"log"
// 	"strconv"
// )

// func main() {
// 	forever := make(chan bool)

// 	//创建一个新的RabbitMQ实例
// 	rabbitmq, err := utils.NewRabbitMQ("queue1", "", "", "amqp://guest:guest@localhost:5672/")
// 	defer rabbitmq.Destroy()
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			rabbitmq.Publish("消息:" + strconv.Itoa(i))
// 		}
// 	}()

// 	<-forever
// }
