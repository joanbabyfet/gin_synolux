package dto

// ==================== Query ====================
type FeedbackQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Content string `json:"content"`
	Limit    int    `form:"limit"`
	Count    bool   `form:"count"`
	IsAdmin	 bool   `form:"is_admin"` // 是否后台请求
}
