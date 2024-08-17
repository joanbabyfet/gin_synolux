package dto

import "gin-synolux/utils"

type MovieQuery struct {
	utils.Pager
	Title  string `json:"title"`
	Status int    `json:"status"`
}
