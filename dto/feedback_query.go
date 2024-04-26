package dto

import "gin-synolux/utils"

type FeedbackQuery struct {
	utils.Pager
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Content string `json:"content"`
}
