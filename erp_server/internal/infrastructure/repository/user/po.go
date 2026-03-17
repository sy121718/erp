package repository

import (
	"time"

	"erp-server/internal/domain/user"
)

// UserPO 用户持久化对象
type UserPO struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Username      string     `gorm:"column:username;uniqueIndex"`
	Password      string     `gorm:"column:password"`
	Name          string     `gorm:"column:name"`
	Avatar        string     `gorm:"column:avatar"`
	Email         string     `gorm:"column:email"`
	Phone         string     `gorm:"column:phone"`
	Status        int        `gorm:"column:status"`
	Points        int64      `gorm:"column:points"`
	TotalPoints   int64      `gorm:"column:total_points"`
	UsedPoints    int64      `gorm:"column:used_points"`
	ExpireTime    *time.Time `gorm:"column:expire_time"`
	LastLoginTime *time.Time `gorm:"column:last_login_time"`
	LastLoginIP   string     `gorm:"column:last_login_ip"`
	CreateBy      int64      `gorm:"column:create_by"`
	CreateTime    time.Time  `gorm:"column:create_time;autoCreateTime:milli"`
	UpdateBy      int64      `gorm:"column:update_by"`
	UpdateTime    time.Time  `gorm:"column:update_time;autoUpdateTime:milli"`
	DeletedTime   *time.Time `gorm:"column:deleted_time"`
}

// TableName 指定表名
func (UserPO) TableName() string {
	return "user"
}

// ToEntity 转换为实体
func (po *UserPO) ToEntity() *user.User {
	return &user.User{
		ID:            po.ID,
		Username:      po.Username,
		Password:      po.Password,
		Name:          po.Name,
		Avatar:        po.Avatar,
		Email:         po.Email,
		Phone:         po.Phone,
		Status:        user.Status(po.Status),
		Points:        po.Points,
		TotalPoints:   po.TotalPoints,
		UsedPoints:    po.UsedPoints,
		ExpireTime:    po.ExpireTime,
		LastLoginTime: po.LastLoginTime,
		LastLoginIP:   po.LastLoginIP,
		CreateBy:      po.CreateBy,
		CreateTime:    po.CreateTime,
		UpdateBy:      po.UpdateBy,
		UpdateTime:    po.UpdateTime,
	}
}

// ToPO 实体转换为PO
func ToPO(entity *user.User) *UserPO {
	return &UserPO{
		ID:            entity.ID,
		Username:      entity.Username,
		Password:      entity.Password,
		Name:          entity.Name,
		Avatar:        entity.Avatar,
		Email:         entity.Email,
		Phone:         entity.Phone,
		Status:        int(entity.Status),
		Points:        entity.Points,
		TotalPoints:   entity.TotalPoints,
		UsedPoints:    entity.UsedPoints,
		ExpireTime:    entity.ExpireTime,
		LastLoginTime: entity.LastLoginTime,
		LastLoginIP:   entity.LastLoginIP,
		CreateBy:      entity.CreateBy,
		UpdateBy:      entity.UpdateBy,
	}
}
