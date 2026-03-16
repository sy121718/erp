package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Init 初始化验证器
func Init() {
	validate = validator.New()
}

// Get 获取验证器
func Get() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	return validate
}

// Validate 验证结构体
func Validate(s interface{}) error {
	return Get().Struct(s)
}

// ValidateVar 验证变量
func ValidateVar(field interface{}, tag string) error {
	return Get().Var(field, tag)
}
