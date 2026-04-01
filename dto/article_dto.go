package dto

// ==================== Query ====================
type ArticleQuery struct {
	Status   *int8  `form:"status"`
	Catid    int    `form:"catid"`
	Catids   []int  `form:"catids[]"`
	Title    string `form:"title"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Limit    int    `form:"limit"`
	Count    bool   `form:"count"`
	IsAdmin	 bool   `form:"is_admin"` // 是否后台请求
}

// ==================== Response ====================
type ArticleDeleteReq struct {
	ID     int
	UserID string
	Role   string
}

type ArticleChangeStatusReq struct {
    ID      int 
    Status  int 
	UserID string
	Role   string
}

// ==================== Response ====================
type ArticleListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}