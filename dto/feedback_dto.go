package dto

// ==================== Query ====================
type FeedbackQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Limit    int    `form:"limit,default=0"`
	Name     string `form:"name"`
	Mobile   string `form:"mobile"`
	Email    string `form:"email"`
	Content  string `form:"content"`

	// 内部字段
	Count bool `form:"-"`
}

// ==================== Request ====================
type FeedbackSaveReq struct {
	Name    string `form:"name" binding:"required"`
	Mobile  string `form:"mobile" binding:"required"`
	Email   string `form:"email" binding:"required"`
	Content string `form:"content" binding:"required"`
	CreateUser string	`form:"create_user"`
}
