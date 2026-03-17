package repository

import (
	"context"
	"errors"

	"erp-server/internal/domain/admin"
	"erp-server/pkg/database"

	"gorm.io/gorm"
)

// AdminRepository 管理员仓储实现
type AdminRepository struct{}

// NewAdminRepository 创建管理员仓储
func NewAdminRepository() admin.AdminRepository {
	return &AdminRepository{}
}

// getDB 获取数据库连接（支持自动重连）
func (r *AdminRepository) getDB() *gorm.DB {
	return database.Get()
}

// FindByID 根据ID查询
func (r *AdminRepository) FindByID(ctx context.Context, id int64) (*admin.Admin, error) {
	var po AdminPO
	err := r.getDB().WithContext(ctx).Where("id = ?", id).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToEntity(), nil
}

// FindByUsername 根据用户名查询
func (r *AdminRepository) FindByUsername(ctx context.Context, username string) (*admin.Admin, error) {
	var po AdminPO
	err := r.getDB().WithContext(ctx).Where("username = ?", username).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToEntity(), nil
}

// List 查询列表（过滤掉超管）
func (r *AdminRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]*admin.Admin, int64, error) {
	var pos []AdminPO
	var total int64

	query := r.getDB().WithContext(ctx).Model(&AdminPO{})

	// 过滤掉超管（is_admin = false）
	query = query.Where("is_admin = ?", false)

	if keyword != "" {
		query = query.Where("username LIKE ? OR name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	// 转换为实体
	entities := make([]*admin.Admin, len(pos))
	for i, po := range pos {
		entities[i] = po.ToEntity()
	}

	return entities, total, nil
}

// Save 保存
func (r *AdminRepository) Save(ctx context.Context, entity *admin.Admin) error {
	po := ToPO(entity)
	return r.getDB().WithContext(ctx).Create(po).Error
}

// Update 更新
func (r *AdminRepository) Update(ctx context.Context, entity *admin.Admin) error {
	po := ToPO(entity)
	return r.getDB().WithContext(ctx).Save(po).Error
}

// Delete 删除
func (r *AdminRepository) Delete(ctx context.Context, id int64) error {
	return r.getDB().WithContext(ctx).Delete(&AdminPO{}, id).Error
}

// UpdatePassword 更新密码
func (r *AdminRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	return r.getDB().WithContext(ctx).Model(&AdminPO{}).Where("id = ?", id).Update("password", password).Error
}

// UpdateStatus 更新状态
func (r *AdminRepository) UpdateStatus(ctx context.Context, id int64, status admin.Status) error {
	return r.getDB().WithContext(ctx).Model(&AdminPO{}).Where("id = ?", id).Update("status", int(status)).Error
}