package main

import (
	"gin-synolux/models"
	"gin-synolux/routers"
	"gin-synolux/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	// 设置gin运行模式
	gin.SetMode(gin.DebugMode)

	//数据库初始化
	models.DB.Init()
	defer models.DB.Close() //延迟关闭

	//初始化定时任务
	//utils.CrontabInit()

	//初始化工作队列
	utils.InitRedisQueue()

	//路由
	router := routers.Init()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
