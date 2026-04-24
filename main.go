package main

import (
	"gin-synolux/common"
	"gin-synolux/db"
	"gin-synolux/routers"
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
	if err := common.InitConfig(*cfg); err != nil {
		panic(err)
	}

	// 这里初始化 logger（最佳位置）
	common.InitLogger()

	// 设置gin运行模式
	gin.SetMode(gin.DebugMode)

	//初始化数据库
	if err := db.DB.Init(); err != nil {
		panic("数据库未初始化成功")
	}

	defer db.DB.Close() //延迟关闭

	//初始化定时任务
	//common.CrontabInit()

	//初始化工作队列
	//common.InitRedisQueue()

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
