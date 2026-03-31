package dto

// ==================== Query ====================
type AdQuery struct {
	Title    string `form:"title"`
	Catid   int   `json:"catid"`
	Type    int   `json:"type"`
	Status  *int   `json:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Limit   int   `json:"limit"`
	Catids  []int `json:"catids"`
	Count   bool  `json:"count"`
	IsAdmin bool  `json:"is_admin"` // 是否后台请求
}

// ==================== Response ====================
type AdListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}
