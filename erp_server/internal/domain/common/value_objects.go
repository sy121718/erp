package common

import (
	"time"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status 状态枚举
type Status int

const (
	StatusDisabled Status = 0
	StatusEnabled  Status = 1
)

// IsActive 是否启用
func (s Status) IsActive() bool {
	return s == StatusEnabled
}
