package dto

import "gin-synolux/utils"

type AdQuery struct {
	utils.Pager
	Catid  int `json:"catid"`
	Type   int `json:"type"`
	Status int `json:"status"`
	Limit  int    `json:"limit"`
	Catids []int  `json:"catids"`
	Count	bool  `json:"count"`
	IsAdmin bool   `json:"is_admin"` // 是否后台请求
}
