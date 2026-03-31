package common

import (
	"errors"

	"github.com/thedevsaddam/govalidator"
)

// ValidateStruct 封装 govalidator 验证逻辑
// data: 需要验证的结构体指针
// rules: 验证规则 map[string][]string
// messages: 验证错误消息 map[string][]string
// 返回第一个错误，如果全部通过返回 nil
func ValidateStruct(data interface{}, rules, messages govalidator.MapData) error {
    opts := govalidator.Options{
        Data:            data,
        Rules:           rules,
        Messages:        messages,
        RequiredDefault: false,
    }
    valid := govalidator.New(opts)
    errs := valid.ValidateStruct()
    if len(errs) > 0 {
        for _, v := range errs {
            return errors.New(v[0])
        }
    }
    return nil
}