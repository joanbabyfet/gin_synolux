package dto

// ==================== Query ====================
type AdminQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username" description:"账号"` //账号
	Status   *int    `form:"status" description:"状态"`    //状态
	Limit    int    `form:"limit"`
	Count    bool   `form:"count"`
	IsAdmin	 bool   `form:"is_admin"` // 是否后台请求
}

// ==================== Request ====================
type AdminDetailReq struct {
	UID string `form:"id" binding:"required"`
}

type AdminLoginReq struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code"` //验证码
	Key      string `form:"key"`	//验证码key
	LoginIp      string `form:"ip"`
}

type AdminSetPasswordReq struct {
	Password    string `form:"password" binding:"required"`		//原始密码
	NewPassword string `form:"new_password" binding:"required"`	//新密码
	RePassword  string `form:"re_password" binding:"required"`	//确认密码
	UID    		string `form:"id"`								//用户id
}

type AdminRegisterReq struct {
	Username  string `form:"username"`
	Password  string `form:"password"`
	Realname  string `form:"realname"`
	Email     string `form:"email"`
	Sex       int8   `form:"sex"`
	RegIp 	string	`form:"reg_ip"`
}

type AdminProfileReq struct {
	ID       string `form:"id"` // 有值=更新，无值=新增
	Password  string `form:"password"`
	Realname  string `form:"realname"`
	Email     string `form:"email"`
	Sex       int8   `form:"sex"`
}

// ==================== Response ====================
type AdminListResp struct {
	List  interface{} `json:"list"`
	Count int64       `json:"count"`
}

type AdminLoginResp struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Realname  string `json:"realname"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}