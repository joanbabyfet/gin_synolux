package db

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Database struct {
	Self *gorm.DB
}

//定义一个全局变量, 创建一个 Database 实例并获取内存地址
var DB = &Database{}
//单例 + 只初始化一次
var once sync.Once

// 初始化数据库
func (db *Database) Init() error {
	var err error

	once.Do(func() {
		db.Self, err = openDB(
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.addr"),
			viper.GetString("db.name"),
			viper.GetString("db.charset"),
		)
	})
	return err
}

func openDB(username, password, addr, name, charset string) (*gorm.DB, error) {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		charset,
		true,
		"Local")

	db, err := gorm.Open("mysql", config)

	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
		return nil, err
	}
	setupDB(db)

	//检查连接(Ping 检测)
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// 配置连接池
func setupDB(db *gorm.DB) {
	if db == nil {
		return
	}

	db.LogMode(viper.GetBool("gormlog"))
	//db.LogMode(true)
	db.DB().SetMaxIdleConns(10) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
	db.DB().SetMaxOpenConns(100) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误
}

// 关闭
func (db *Database) Close() {
	if db.Self != nil {
		db.Self.Close()
	}
}