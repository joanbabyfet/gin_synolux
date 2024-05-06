package jobs

type HandlerFunc func(args interface{}) error

func (f HandlerFunc) CallFunc(args interface{}) error {
	return f(args)
}
