package repository

import (
	"time"

	"erp-server/internal/domain/admin"
)

// AdminPO 管理员持久化对象
type AdminPO struct {
	ID                int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Username          string         `gorm:"column:username;uniqueIndex"`
	Password          string         `gorm:"column:password"`
	Name              string         `gorm:"column:name"`
	Avatar            string         `gorm:"column:avatar"`
	Email             string         `gorm:"column:email"`
	Phone             string         `gorm:"column:phone"`
	Status            int            `gorm:"column:status"`
	IsAdmin           int            `gorm:"column:is_admin"`
	LoginFailureCount int            `gorm:"column:login_failure_count"`
	LockedUntilTime   *time.Time     `gorm:"column:locked_until_time"`
	LastFailureTime   *time.Time     `gorm:"column:last_failure_time"`
	RegisterIP        string         `gorm:"column:register_ip"`
	RegisterLocation  string         `gorm:"column:register_location"`
	LastLoginIP       string         `gorm:"column:last_login_ip"`
	LastLoginLocation string         `gorm:"column:last_login_location"`
	LastLoginISP      string         `gorm:"column:last_login_isp"`
	LastLoginTime     *time.Time     `gorm:"column:last_login_time"`
	CreateBy          int64          `gorm:"column:create_by"`
	CreateTime        time.Time      `gorm:"column:create_time;autoCreateTime:milli"`
	UpdateBy          int64          `gorm:"column:update_by"`
	UpdateTime        time.Time      `gorm:"column:update_time;autoUpdateTime:milli"`
	DeletedTime       *time.Time     `gorm:"column:deleted_time"`
}

// TableName 指定表名
func (AdminPO) TableName() string {
	return "sys_admin"
}

// ToEntity 转换为实体
func (po *AdminPO) ToEntity() *admin.Admin {
	return &admin.Admin{
		ID:                po.ID,
		Username:          po.Username,
		Password:          po.Password,
		Name:              po.Name,
		Avatar:            po.Avatar,
		Email:             po.Email,
		Phone:             po.Phone,
		Status:            admin.Status(po.Status),
		IsAdmin:           po.IsAdmin == 1,
		LoginFailureCount: po.LoginFailureCount,
		LockedUntilTime:   po.LockedUntilTime,
		LastFailureTime:   po.LastFailureTime,
		RegisterIP:        po.RegisterIP,
		RegisterLocation:  po.RegisterLocation,
		LastLoginIP:       po.LastLoginIP,
		LastLoginLocation: po.LastLoginLocation,
		LastLoginISP:      po.LastLoginISP,
		LastLoginTime:     po.LastLoginTime,
		CreateBy:          po.CreateBy,
		CreateTime:        po.CreateTime,
		UpdateBy:          po.UpdateBy,
		UpdateTime:        po.UpdateTime,
	}
}

// ToPO 实体转换为PO
func ToPO(entity *admin.Admin) *AdminPO {
	isAdmin := 0
	if entity.IsAdmin {
		isAdmin = 1
	}
	return &AdminPO{
		ID:                entity.ID,
		Username:          entity.Username,
		Password:          entity.Password,
		Name:              entity.Name,
		Avatar:            entity.Avatar,
		Email:             entity.Email,
		Phone:             entity.Phone,
		Status:            int(entity.Status),
		IsAdmin:           isAdmin,
		LoginFailureCount: entity.LoginFailureCount,
		LockedUntilTime:   entity.LockedUntilTime,
		LastFailureTime:   entity.LastFailureTime,
		RegisterIP:        entity.RegisterIP,
		RegisterLocation:  entity.RegisterLocation,
		LastLoginIP:       entity.LastLoginIP,
		LastLoginLocation: entity.LastLoginLocation,
		LastLoginISP:      entity.LastLoginISP,
		LastLoginTime:     entity.LastLoginTime,
		CreateBy:          entity.CreateBy,
		UpdateBy:          entity.UpdateBy,
	}
}