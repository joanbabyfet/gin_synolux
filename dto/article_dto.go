package dto

// ==================== Query ====================
type ArticleQuery struct {
	Title    string `form:"title"`
	Catid    int    `form:"catid,default=0"`
	Catids   []int  `form:"catids[]"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Limit    int    `form:"limit,default=0"`
	// 内部控制
	Status *int8 `json:"-"`
	Count  bool  `json:"-"` // 是否返回总数
}

// ==================== Request ====================
type ArticleDetailReq struct {
	ID int `form:"id" binding:"required"`
}

type ArticleSaveReq struct {
	ID      int    `form:"id"`
	Catid   int    `form:"catid"`
	Title   string `form:"title" binding:"required"`
	Info    string `form:"info"`
	Content string `form:"content"`
	Author  string `form:"author"`
	Status  int    `form:"status"`
	CreateUser string	`form:"create_user"`
	UpdateUser string	`form:"update_user"`
}

type ArticleDeleteReq struct {
	ID int `json:"id" binding:"required"`

	UserID string `json:"-"`
	Role   string `json:"-"`
}

type ArticleChangeStatusReq struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required"`

	UserID string `json:"-"`
	Role   string `json:"-"`
}

// ==================== Response ====================
type ArticleListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}