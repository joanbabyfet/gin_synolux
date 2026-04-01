package dto

// ==================== Query ====================
type UserQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `json:"username" description:"账号"` //账号
	Status   *int    `json:"status" description:"状态"`    //状态
	Limit    int    `form:"limit"`
	Count    bool   `form:"count"`
	IsAdmin	 bool   `form:"is_admin"` // 是否后台请求
}

// ==================== Request ====================
type UserLoginReq struct {
	Username string `json:"username" description:"账号"`
	Password string `json:"password" description:"密码"`
	Code     string `json:"code" description:"验证码"`
	Key      string `json:"key" description:"验证码key"`
	LoginIp  string `json:"login_ip" description:"最后登录IP"`
}

type UserPasswordReq struct {
	Password    string `json:"password"`     //原始密码
	NewPassword string `json:"new_password"` //新密码
	RePassword  string `json:"re_password"`  //确认密码
	Uid         string `json:"uid"`          //用户id
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