package service

type ServiceError struct {
	Code int
	Msg  string
}

func (e *ServiceError) Error() string {
	return e.Msg
}

func NewServiceError(code int, msg string) *ServiceError {
	return &ServiceError{
		Code: code,
		Msg:  msg,
	}
}