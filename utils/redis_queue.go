package utils

import (
	"encoding/json"
	"strconv"

	"github.com/lexkong/log"
)

// 工作任务的结构体
type Job struct {
	Queue string                 `json:"queue"`
	Task  string                 `json:"task"` //任务名, 唯一标识
	Args  map[string]interface{} `json:"args"` //参数
}

// 初始化工作队列
func InitRedisQueue() {
	// 初始化工作队列
	cache_key := "work_queue"
	_, err := Redis.Del(cache_key).Result()
	if err != nil {
		panic(err)
	}

	// 启动多个worker，同时监听工作队列
	for i := 1; i <= 5; i++ {
		//走协程
		go worker(i)
	}
}

// 消费者, 工作处理函数
func worker(id int) {
	for {
		// 从工作队列中获取一个工作
		cache_key := "work_queue"                       //队列名
		item, err := Redis.BLPop(0, cache_key).Result() //移出并获取列表的第一个元素
		if err != nil {
			log.Error("worker "+strconv.Itoa(id)+" failed", err)
			continue
		}

		// 解码工作的JSON数据
		var job Job
		jobJSON := item[1]                          //这里只获取值不获取key
		err = json.Unmarshal([]byte(jobJSON), &job) //json字符串转struct
		if err != nil {
			log.Error("worker "+strconv.Itoa(id)+" failed", err)
			continue
		}

		//发送邮件
		if job.Task == "send_mail" {
			to := job.Args["to"].(string)
			subject := job.Args["subject"].(string)
			body := job.Args["body"].(string)
			ok := SendMail(to, subject, body)
			if !ok {
				log.Error("worker "+strconv.Itoa(id)+" process: "+job.Task+" failed args: "+jobJSON, err)
				continue
			}
		}
		if job.Task == "send_sms" {
			to := job.Args["to"].(string)
			body := job.Args["body"].(string)
			ok := SendSMS(to, body)
			if !ok {
				log.Error("worker "+strconv.Itoa(id)+" process: "+job.Task+" failed args: "+jobJSON, err)
				continue
			}
		}
		log.Info("Worker " + strconv.Itoa(id) + ": Processing job " + jobJSON + ")\n")
	}
}
