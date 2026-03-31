package models

import (
	"github.com/spf13/viper"
)

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type Feedback struct {
	Id         int    `gorm:"primary_key;auto_increment;default();description(ID)" json:"id"`
	Name       string `gorm:"size(50);default();null;index;description(姓名)" json:"name"`
	Mobile     string `gorm:"size(20);default();null;index;description(手机号)" json:"mobile"`
	Email      string `gorm:"size(100);default();null;index;description(信箱)" json:"email"`
	Content    string `gorm:"type(text);null;description(內容)" json:"content"`
	CreateTime int    `gorm:"default(0);null;description(創建時間)" json:"create_time"`
	CreateUser string `gorm:"size(32);default(0);null;description(創建人)" json:"create_user"`
	UpdateTime int    `gorm:"default(0);null;description(修改時間)" json:"update_time"`
	UpdateUser string `gorm:"size(32);default(0);null;description(修改人)" json:"update_user"`
	DeleteTime int    `gorm:"default(0);null;description(刪除時間)" json:"delete_time"`
	DeleteUser string `gorm:"size(32);default(0);null;description(刪除人)" json:"delete_user"`
}

func (m *Feedback) TableName() string {
	return viper.GetString("db.prefix") + "feedback"
}