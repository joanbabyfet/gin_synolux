package common

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