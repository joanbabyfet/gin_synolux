package dto

import "gin-synolux/utils"

type AdQuery struct {
	utils.Pager
	Catid  int `json:"catid"`
	Type   int `json:"type"`
	Status int `json:"status"`
}
