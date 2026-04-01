package dto

// ==================== Query ====================
type MovieQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Title  string `json:"title"`
	Status *int8    `json:"status"`
	Limit  int    `json:"limit"`
	Count	bool  `json:"count"`
	IsAdmin bool   `json:"is_admin"` // 是否后台请求
}

// ==================== Response ====================
type MovieDeleteReq struct {
	ID     int
	UserID string
	Role   string
}

type MovieChangeStatusReq struct {
    ID      int 
    Status  int 
	UserID string
	Role   string
}

// ==================== Response ====================
type MovieListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}
