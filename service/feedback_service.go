package service

import (
	"errors"
	"gin-synolux/models"
	"gin-synolux/utils"

	"github.com/lexkong/log"
)

type FeedbackService struct {
}

// 保存
func (s *FeedbackService) Save(data models.Feedback) (int, error) {
	stat := 1

	if data.Id > 0 {

	} else {
		data.CreateUser = "1"               //添加人
		data.CreateTime = utils.Timestamp() //添加时间
		err := data.Add()
		if err != nil {
			log.Error("反馈添加失败", err)
			return -2, errors.New("反馈添加失败")
		}
	}
	return stat, nil
}
