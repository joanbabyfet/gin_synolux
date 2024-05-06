package queue

import (
	"bytes"
	"encoding/gob"
	"gin-synolux/jobs"
	"log"

	"github.com/streadway/amqp"
)

type Queue struct {
	QueueConn *amqp.Connection
	JobPool   map[string]jobs.HandlerFunc
}

func NewQueue() *Queue {
	q := new(Queue)
	q.Connect()
	q.JobPool = make(map[string]jobs.HandlerFunc)
	return q
}

func (q *Queue) Connect() {
	var err error

	url := "amqp://guest:guest@localhost:5672/"
	q.QueueConn, err = amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
}

func (q *Queue) PushJob(job string, f jobs.HandlerFunc) {
	q.JobPool[job] = f
}

func (q *Queue) Close() {
	q.QueueConn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

/**
 * 多任务共享一个队列
 */
func (queue *Queue) NewShareQueue(queueName string) {
	log.Printf(" [*] Waiting for %s messages. To exit press CTRL+C", queueName)

	ch, err := queue.QueueConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		//接收消息
		var sender Sender
		decoder := gob.NewDecoder(bytes.NewReader(msg.Body))
		decoder.Decode(&sender)
		//fmt.Println(string(sender.QueueName), string(sender.Job),string(sender.Msg))

		//查job处理方法
		if f, ok := queue.JobPool[string(sender.Job)]; ok {
			f(sender.Msg)
			msg.Ack(false)
		}
	}
}
