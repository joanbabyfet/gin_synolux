package dto

// 登录请求格式
type UserLogin struct {
	Username string `json:"username" description:"账号"`
	Password string `json:"password" description:"密码"`
	Code     string `json:"code" description:"验证码"`
	Key      string `json:"key" description:"验证码key"`
	LoginIp  string `json:"login_ip" description:"最后登录IP"`
}

// 登录响应格式
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
