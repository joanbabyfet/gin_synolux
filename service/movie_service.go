package service

import (
	"errors"
	"gin-synolux/dto"
	"gin-synolux/models"
	"strconv"

	"github.com/lexkong/log"
)

type MovieService struct {
}

// 获取分页列表
func (s *MovieService) List(query dto.MovieQuery) ([]*models.Movie, int64) {
	entity := new(models.Movie) //new实例化
	return entity.List(query)
}

// 获取详情
func (s *MovieService) GetById(id int) (*models.Movie, error) {
	entity := new(models.Movie)
	info, err := entity.GetById(id)
	if err != nil {
		log.Error("视频不存在 "+strconv.Itoa(id), err)
		return nil, errors.New("视频不存在")
	}
	return info, nil
}
