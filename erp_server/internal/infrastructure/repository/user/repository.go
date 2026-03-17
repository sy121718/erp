package repository

import (
	"context"
	"errors"

	"erp-server/internal/domain/user"
	"erp-server/pkg/database"

	"gorm.io/gorm"
)

// UserRepository 用户仓储实现
type UserRepository struct{}

// NewUserRepository 创建用户仓储
func NewUserRepository() user.UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) getDB() *gorm.DB {
	return database.Get()
}

// FindByID 根据ID查询
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var po UserPO
	err := r.getDB().WithContext(ctx).Where("id = ? AND deleted_time IS NULL", id).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToEntity(), nil
}

// FindByUsername 根据用户名查询
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var po UserPO
	err := r.getDB().WithContext(ctx).Where("username = ? AND deleted_time IS NULL", username).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToEntity(), nil
}

// List 查询列表
func (r *UserRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]*user.User, int64, error) {
	var pos []UserPO
	var total int64

	query := r.getDB().WithContext(ctx).Model(&UserPO{}).Where("deleted_time IS NULL")

	if keyword != "" {
		query = query.Where("username LIKE ? OR name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*user.User, len(pos))
	for i, po := range pos {
		entities[i] = po.ToEntity()
	}

	return entities, total, nil
}

// Save 保存
func (r *UserRepository) Save(ctx context.Context, entity *user.User) error {
	po := ToPO(entity)
	result := r.getDB().WithContext(ctx).Create(po)
	if result.Error != nil {
		return result.Error
	}
	entity.ID = po.ID
	return nil
}

// Update 更新
func (r *UserRepository) Update(ctx context.Context, entity *user.User) error {
	po := ToPO(entity)
	return r.getDB().WithContext(ctx).Save(po).Error
}

// Delete 软删除
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	now := gorm.Expr("NOW(3)")
	return r.getDB().WithContext(ctx).Model(&UserPO{}).Where("id = ?", id).Update("deleted_time", now).Error
}

// UpdatePassword 更新密码
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	return r.getDB().WithContext(ctx).Model(&UserPO{}).Where("id = ?", id).Update("password", password).Error
}

// UpdateStatus 更新状态
func (r *UserRepository) UpdateStatus(ctx context.Context, id int64, status user.Status) error {
	return r.getDB().WithContext(ctx).Model(&UserPO{}).Where("id = ?", id).Update("status", int(status)).Error
}
