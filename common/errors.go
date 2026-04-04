package common

import "github.com/go-playground/validator/v10"

type ServiceError struct {
	Code int
	Msg  string
	Err  error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Msg
}

func NewError(code int, msg string) *ServiceError {
	return &ServiceError{
		Code: code,
		Msg:  msg,
	}
}

func WrapError(err error, code int, msg string) *ServiceError {
	return &ServiceError{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}

//gin自带参数校验
func GetValidMsg(err error) string {
    if errs, ok := err.(validator.ValidationErrors); ok {
        e := errs[0]

        switch e.Tag() {
        case "required":
            return e.Field() + "不能为空"
        case "min":
            return e.Field() + "长度太短"
        case "max":
            return e.Field() + "长度太长"
        case "oneof":
            return e.Field() + "值非法"
        }
    }
    return err.Error()
}