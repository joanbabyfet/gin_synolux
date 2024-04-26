package dto

type Password struct {
	Password    string `json:"password"`     //原始密码
	NewPassword string `json:"new_password"` //新密码
	RePassword  string `json:"re_password"`  //确认密码
	Uid         string `json:"uid"`          //用户id
}
