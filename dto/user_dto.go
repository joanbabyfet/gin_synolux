package dto

// ==================== Query ====================
type UserQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username" description:"账号"` //账号
	Status   *int    `form:"status" description:"状态"`    //状态
	Limit    int    `form:"limit"`
	Count    bool   `form:"count"`
	IsAdmin	 bool   `form:"is_admin"` // 是否后台请求
}

// ==================== Request ====================
type UserDetailReq struct {
	UID string `form:"id" binding:"required"`
}

type UserLoginReq struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code"` //验证码
	Key      string `form:"key"`	//验证码key
	Ip      string `form:"ip"`
}

type UserSetPasswordReq struct {
	Password    string `form:"password" binding:"required"`
	NewPassword string `form:"new_password" binding:"required"`
	RePassword  string `form:"re_password" binding:"required"`
	UID    		string `form:"id"`
}

type UserRegisterReq struct {
	Username  string `form:"username"`
	Password  string `form:"password"`
	Realname  string `form:"realname"`
	Email     string `form:"email"`
	PhoneCode string `form:"phone_code"`
	Phone     string `form:"phone"`
	Avatar    string `form:"avatar"`
	Sex       int8   `form:"sex"`
	RegIp 	string	`form:"reg_ip"`
}

type UserProfileReq struct {
	ID       string `form:"id"` // 有值=更新，无值=新增
	Password  string `form:"password"`
	Realname  string `form:"realname"`
	Email     string `form:"email"`
	PhoneCode string `form:"phone_code"`
	Phone     string `form:"phone"`
	Avatar    string `form:"avatar"`
	Sex       int8   `form:"sex"`
}

// ==================== Response ====================
type UserListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}

type UserLoginResp struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Realname  string `json:"realname"`
	Email     string `json:"email"`
	PhoneCode string `json:"phone_code"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	Language  string `json:"language"`
	Token     string `json:"token"`
}