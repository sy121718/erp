package user

import "time"

// Status 用户状态
type Status int

const (
	StatusDisabled Status = 0 // 禁用
	StatusEnabled  Status = 1 // 启用
)

// User 用户实体
type User struct {
	ID            int64
	Username      string
	Password      string
	Name          string
	Avatar        string
	Email         string
	Phone         string
	Status        Status
	Points        int64
	TotalPoints   int64
	UsedPoints    int64
	ExpireTime    *time.Time
	LastLoginTime *time.Time
	LastLoginIP   string
	CreateBy      int64
	CreateTime    time.Time
	UpdateBy      int64
	UpdateTime    time.Time
}

// IsActive 是否激活
func (u *User) IsActive() bool {
	return u.Status == StatusEnabled
}

// CanLogin 是否可以登录
func (u *User) CanLogin() bool {
	return u.IsActive()
}

// UpdateLastLogin 更新最后登录信息
func (u *User) UpdateLastLogin(ip string) {
	u.LastLoginIP = ip
	u.LastLoginTime = ptrTime(time.Now())
}

// Disable 禁用账户
func (u *User) Disable() {
	u.Status = StatusDisabled
}

// Enable 启用账户
func (u *User) Enable() {
	u.Status = StatusEnabled
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
