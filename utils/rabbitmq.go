package utils

import (
	"errors"

	"github.com/lexkong/log"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn         *amqp.Connection //连接
	channel      *amqp.Channel    //管道
	QueueName    string           //队列名称
	ExchangeName string           //交换机名称
	key          string           //Binding Key/Routing Key, Simple模式 几乎用不到
	MqURL        string           //连接信息-amqp://账号:密码@地址:端口号/-amqp://guest:guest@127.0.0.1:5672/
}

// 创建一个新的RabbitMQ实例
func NewRabbitMQ(queueName, exchangeName, key, mqurl string) (*RabbitMQ, error) {
	var err error
	if queueName == "" || mqurl == "" {
		return nil, errors.New("QueueName and mqUrl is required")
	}

	r := &RabbitMQ{
		QueueName:    queueName,
		ExchangeName: exchangeName,
		key:          key,
		MqURL:        mqurl,
	}

	// 创建连接
	r.conn, err = amqp.Dial(r.MqURL)
	if err != nil {
		log.Error("failed to connect to RabbitMQ", err)
		return nil, err
	}
	// 创建channel
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Error("failed to open a channel", err)
		return nil, err
	}
	return r, nil
}

// 简单模式：发送消息
func (r *RabbitMQ) Publish(message string) error {
	// 声明队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("failed to declare a queue", err)
		return err
	}
	// 发送消息到队列中
	err = r.channel.Publish(
		r.ExchangeName,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain", //没有声明默认为text/plain
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Error("failed to publish a message", err)
		return err
	}
	return nil
}

// 简单模式：接收消息
func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	// 声明队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,   //额外的属性
	)
	if err != nil {
		log.Error("failed to declare a queue", err)
		return nil, err
	}
	// 消费消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("failed to register a consumer", err)
		return nil, err
	}
	return msgs, nil
}

// 错误处理
// func (r *RabbitMQ) failOnErr(err error, msg string) {
// 	if err != nil {
// 		log.Error(msg, err)
// 		panic(err)
// 	}
// }

// 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
	log.Infof("%s,%s is closed!!!", r.ExchangeName, r.QueueName)
}
