package dto

import "gin-synolux/utils"

type MovieQuery struct {
	utils.Pager
	Title  string `json:"title"`
	Status int    `json:"status"`
	Limit  int    `json:"limit"`
	Count	bool  `json:"count"`
	IsAdmin bool   `json:"is_admin"` // 是否后台请求
}
