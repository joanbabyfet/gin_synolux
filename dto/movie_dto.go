package dto

// ==================== Query ====================
type MovieQuery struct {
	Title    string `form:"title"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Limit    int    `form:"limit,default=0"`

	// 内部控制
	Status *int8 `json:"-"`
	Count  bool  `json:"-"`
}

// ==================== Request ====================
type MovieDeleteReq struct {
	ID int `json:"id" binding:"required"`

	UserID string `json:"-"`
	Role   string `json:"-"`
}

type MovieDetailReq struct {
	ID int `form:"id" binding:"required"`
}

type MovieSaveReq struct {
	ID     int    `form:"id"`
	Title  string `form:"title" binding:"required"`
	Img    string `form:"img"`
	Url    string `form:"url"`
	Sort   int    `form:"sort"`
	Status int    `form:"status"`
	CreateUser string	`form:"create_user"`
	UpdateUser string	`form:"update_user"`
}

type MovieChangeStatusReq struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required"`

	UserID string `json:"-"`
	Role   string `json:"-"`
}

// ==================== Response ====================
type MovieListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}
