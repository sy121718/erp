package admin

import "time"

// Status 管理员状态
type Status int

const (
	StatusDisabled Status = 0 // 禁用
	StatusEnabled  Status = 1 // 启用
)

// Admin 管理员实体
type Admin struct {
	ID                int64
	Username          string
	Password          string
	Name              string
	Avatar            string
	Email             string
	Phone             string
	Status            Status
	IsAdmin           bool
	LoginFailureCount int
	LockedUntilTime   *time.Time
	LastFailureTime   *time.Time
	RegisterIP        string
	RegisterLocation  string
	LastLoginIP       string
	LastLoginLocation string
	LastLoginISP      string
	LastLoginTime     *time.Time
	CreateBy          int64
	CreateTime        time.Time
	UpdateBy          int64
	UpdateTime        time.Time
}

// IsLocked 是否被锁定
func (a *Admin) IsLocked() bool {
	if a.LockedUntilTime == nil {
		return false
	}
	return a.LockedUntilTime.After(time.Now())
}

// IsActive 是否激活
func (a *Admin) IsActive() bool {
	return a.Status == StatusEnabled && !a.IsLocked()
}

// CanLogin 是否可以登录
func (a *Admin) CanLogin() bool {
	return a.IsActive()
}

// IncrementLoginFailure 增加登录失败次数
func (a *Admin) IncrementLoginFailure() {
	a.LoginFailureCount++
	a.LastFailureTime = ptrTime(time.Now())
}

// ResetLoginFailure 重置登录失败次数
func (a *Admin) ResetLoginFailure() {
	a.LoginFailureCount = 0
	a.LastFailureTime = nil
	a.LockedUntilTime = nil
}

// Lock 锁定账户
func (a *Admin) Lock(duration time.Duration) {
	a.LockedUntilTime = ptrTime(time.Now().Add(duration))
}

// Unlock 解锁账户
func (a *Admin) Unlock() {
	a.LockedUntilTime = nil
	a.LoginFailureCount = 0
	a.LastFailureTime = nil
}

// UpdateLastLogin 更新最后登录信息
func (a *Admin) UpdateLastLogin(ip, location, isp string) {
	a.LastLoginIP = ip
	a.LastLoginLocation = location
	a.LastLoginISP = isp
	a.LastLoginTime = ptrTime(time.Now())
}

// Disable 禁用账户
func (a *Admin) Disable() {
	a.Status = StatusDisabled
}

// Enable 启用账户
func (a *Admin) Enable() {
	a.Status = StatusEnabled
}

func ptrTime(t time.Time) *time.Time {
	return &t
}