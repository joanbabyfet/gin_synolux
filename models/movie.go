package models

import (
	"github.com/spf13/viper"
)

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type Movie struct {
	Id         int    `gorm:"primary_key;auto_increment;default();description(ID)" json:"id"`
	Title      string `gorm:"size(50);default();null;index;description(标题)" json:"title"`
	Img        string `gorm:"size(100);default();null;description(图片)" json:"img"`
	Url        string `gorm:"size(100);default();null;description(链接)" json:"url"`
	Sort       int16  `gorm:"default(0);null;description(排序: 数字小的排前面)" json:"sort"`
	Status     int8   `gorm:"default(1);null;description(状态: 0=禁用 1=启用)" json:"status"`
	CreateTime int    `gorm:"default(0);null;description(創建時間)" json:"create_time"`
	CreateUser string `gorm:"size(32);default(0);null;description(創建人)" json:"create_user"`
	UpdateTime int    `gorm:"default(0);null;description(修改時間)" json:"update_time"`
	UpdateUser string `gorm:"size(32);default(0);null;description(修改人)" json:"update_user"`
	DeleteTime int    `gorm:"default(0);null;description(刪除時間)" json:"delete_time"`
	DeleteUser string `gorm:"size(32);default(0);null;description(刪除人)" json:"delete_user"`
}

func (m *Movie) TableName() string {
	return viper.GetString("db.prefix") + "movie"
}
