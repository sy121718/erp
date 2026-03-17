package user

import "context"

// UserRepository 用户仓储接口
type UserRepository interface {
	// FindByID 根据ID查询
	FindByID(ctx context.Context, id int64) (*User, error)

	// FindByUsername 根据用户名查询
	FindByUsername(ctx context.Context, username string) (*User, error)

	// List 查询列表
	List(ctx context.Context, page, pageSize int, keyword string) ([]*User, int64, error)

	// Save 保存
	Save(ctx context.Context, user *User) error

	// Update 更新
	Update(ctx context.Context, user *User) error

	// Delete 删除（软删除）
	Delete(ctx context.Context, id int64) error

	// UpdatePassword 更新密码
	UpdatePassword(ctx context.Context, id int64, password string) error

	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id int64, status Status) error
}
