package dto

// ==================== Query ====================
type AdQuery struct {
	Title    string `form:"title"`
	Catid    int    `form:"catid,default=0"`
	Catids  []int 	`form:"catids"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Limit     int   `form:"limit,default=0"`
	//内部控制
	Status  *int `form:"-"`
	Count   bool `form:"-"` //返回总条数
}

// ==================== Request ====================
type AdDetailReq struct {
	ID int `form:"id" binding:"required"`
}

type AdSaveReq struct {
	ID     int   		 `form:"id"`
	Catid  int    		`form:"catid"`
	Title  string 		`form:"title" binding:"required"`
	Img    string 		`form:"img"`
	Url    string 		`form:"url"`
	Sort   int    		`form:"sort"`
	Status int    		`form:"status"`
	CreateUser string	`form:"create_user"`
	UpdateUser string	`form:"update_user"`
}

type AdDeleteReq struct {
	ID     int    `form:"id" binding:"required"`
	UserID string `form:"-"`
	Role   string `form:"-"`
}

type AdChangeStatusReq struct {
	ID     int    `form:"id" binding:"required"`
	Status int    `form:"-"` //不能让前端传
	UserID string `form:"-"` //不能让前端传, 必须从 JWT 来
	Role   string `form:"-"` //不能让前端传, 必须从 JWT 来
}

// ==================== Response ====================
type AdListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}
