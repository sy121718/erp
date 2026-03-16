package service

import (
	"context"
	"time"

	"erp-server/internal/domain/admin"
	"erp-server/pkg/captcha"
	"erp-server/pkg/errors"
	"erp-server/pkg/jwt"
	"erp-server/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务
type AdminService struct {
	repo       admin.AdminRepository
	tokenStore *token.TokenStore
}

// NewAdminService 创建管理员服务
func NewAdminService(repo admin.AdminRepository) *AdminService {
	return &AdminService{
		repo:       repo,
		tokenStore: token.Get(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username    string
	Password    string
	CaptchaID   string
	CaptchaCode string
	IP          string
	Location    string
	ISP         string
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	Admin        *admin.Admin
}

// Login 登录
func (s *AdminService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 验证验证码
	if req.CaptchaID != "" && req.CaptchaCode != "" {
		if !captcha.Get().Verify(req.CaptchaID, req.CaptchaCode, true) {
			return nil, errors.ErrInvalidParams.WithError(admin.ErrCaptchaInvalid)
		}
	}

	// 查找管理员
	entity, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if entity == nil {
		return nil, errors.ErrUserNotFound.WithError(admin.ErrAdminNotFound)
	}

	// 检查账户状态
	if !entity.CanLogin() {
		if entity.IsLocked() {
			return nil, errors.NewForbidden(403002, "账户已锁定")
		}
		return nil, errors.NewForbidden(403003, "账户已禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password)); err != nil {
		// 增加失败次数
		entity.IncrementLoginFailure()
		// 达到5次锁定30分钟
		if entity.LoginFailureCount >= 5 {
			entity.Lock(30 * time.Minute)
		}
		s.repo.Update(ctx, entity)
		return nil, errors.ErrUnauthorized.WithError(admin.ErrPasswordIncorrect)
	}

	// 重置失败次数并更新登录信息
	entity.ResetLoginFailure()
	entity.UpdateLastLogin(req.IP, req.Location, req.ISP)
	s.repo.Update(ctx, entity)

	// 生成Token
	accessToken, _, err := jwt.Get().GenerateAccessToken(entity.ID, entity.Username, entity.Name, entity.IsAdmin)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	refreshToken, refreshTokenID, err := jwt.Get().GenerateRefreshToken(entity.ID, entity.Username, entity.IsAdmin)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	// 存储refresh_token到Redis
	if err := s.tokenStore.StoreRefreshToken(ctx, entity.ID, refreshTokenID, refreshToken); err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}

	// 存储用户的token_id（用于多设备管理）
	if err := s.tokenStore.StoreUserTokenID(ctx, entity.ID, refreshTokenID); err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin:        entity,
	}, nil
}

// Logout 登出（删除当前用户的refresh_token）
func (s *AdminService) Logout(ctx context.Context, adminID int64, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	// 解析refresh_token获取token_id
	claims, err := jwt.Get().ParseRefreshToken(refreshToken)
	if err != nil {
		return nil // token无效，无需处理
	}

	// 删除Redis中的refresh_token
	_ = s.tokenStore.DeleteRefreshToken(ctx, adminID, claims.TokenID)

	// 从用户token列表中移除
	_ = s.tokenStore.RemoveUserTokenID(ctx, adminID, claims.TokenID)

	return nil
}

// ForceLogout 强制下线（删除用户所有token）
func (s *AdminService) ForceLogout(ctx context.Context, targetID int64) error {
	return s.tokenStore.LogoutUser(ctx, targetID)
}

// GetByID 根据ID获取
func (s *AdminService) GetByID(ctx context.Context, id int64) (*admin.Admin, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if entity == nil {
		return nil, errors.ErrUserNotFound.WithError(admin.ErrAdminNotFound)
	}
	return entity, nil
}

// ListRequest 列表请求
type ListRequest struct {
	Page     int
	PageSize int
	Keyword  string
}

// ListResponse 列表响应
type ListResponse struct {
	List  []*admin.Admin
	Total int64
}

// List 获取列表
func (s *AdminService) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	list, total, err := s.repo.List(ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	return &ListResponse{
		List:  list,
		Total: total,
	}, nil
}

// CreateRequest 创建请求
type CreateRequest struct {
	Username string
	Password string
	Name     string
	Email    string
	Phone    string
	CreateBy int64
}

// Create 创建管理员（is_admin默认为false，不能通过接口创建超管）
func (s *AdminService) Create(ctx context.Context, req *CreateRequest) (*admin.Admin, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	entity := &admin.Admin{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Status:     admin.StatusEnabled,
		IsAdmin:    false, // 默认非超管
		CreateBy:   req.CreateBy,
		UpdateTime: time.Now(),
	}

	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return entity, nil
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	ID       int64
	Name     string
	Email    string
	Phone    string
	UpdateBy int64
}

// Update 更新管理员（只能修改基本信息，不能修改is_admin）
func (s *AdminService) Update(ctx context.Context, req *UpdateRequest) (*admin.Admin, error) {
	entity, err := s.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	entity.Name = req.Name
	entity.Email = req.Email
	entity.Phone = req.Phone
	entity.UpdateBy = req.UpdateBy
	// 注意：不修改 IsAdmin 字段

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return entity, nil
}

// ChangePassword 修改自己的密码
func (s *AdminService) ChangePassword(ctx context.Context, adminID int64, oldPassword, newPassword string) error {
	entity, err := s.GetByID(ctx, adminID)
	if err != nil {
		return err
	}

	// 验证原密码
	if err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(oldPassword)); err != nil {
		return errors.ErrInvalidParams.WithError(admin.ErrOldPasswordIncorrect)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternal.WithError(err)
	}

	if err := s.repo.UpdatePassword(ctx, adminID, string(hashedPassword)); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// ResetPassword 重置管理员密码（仅超管可操作）
func (s *AdminService) ResetPassword(ctx context.Context, operatorID, targetID int64, newPassword string) error {
	// 检查操作者是否是超管
	operator, err := s.GetByID(ctx, operatorID)
	if err != nil {
		return err
	}
	if !operator.IsAdmin {
		return errors.ErrForbidden.WithError(admin.ErrCannotModifyAdmin)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternal.WithError(err)
	}

	if err := s.repo.UpdatePassword(ctx, targetID, string(hashedPassword)); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// Delete 删除管理员
func (s *AdminService) Delete(ctx context.Context, operatorID, targetID int64) error {
	// 不能删除自己
	if operatorID == targetID {
		return errors.ErrForbidden.WithError(admin.ErrCannotDeleteSelf)
	}

	// 检查目标是否存在
	entity, err := s.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	// 不能删除超级管理员
	if entity.IsAdmin {
		return errors.ErrForbidden.WithError(admin.ErrCannotModifyAdmin)
	}

	if err := s.repo.Delete(ctx, targetID); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// Ban 禁用管理员
func (s *AdminService) Ban(ctx context.Context, operatorID, targetID int64) error {
	// 不能禁用自己
	if operatorID == targetID {
		return errors.ErrForbidden.WithError(admin.ErrCannotDisableSelf)
	}

	entity, err := s.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	// 不能禁用超级管理员
	if entity.IsAdmin {
		return errors.ErrForbidden.WithError(admin.ErrCannotModifyAdmin)
	}

	entity.Disable()
	if err := s.repo.Update(ctx, entity); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	// 强制下线该用户
	_ = s.tokenStore.LogoutUser(ctx, targetID)

	return nil
}

// Unban 解禁管理员
func (s *AdminService) Unban(ctx context.Context, targetID int64) error {
	entity, err := s.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	entity.Enable()
	if err := s.repo.Update(ctx, entity); err != nil {
		return errors.ErrDatabase.WithError(err)
	}
	return nil
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

// RefreshToken 刷新Token
func (s *AdminService) RefreshToken(ctx context.Context, oldRefreshToken string) (*RefreshTokenResponse, error) {
	// 解析refresh_token
	claims, err := jwt.Get().ParseRefreshToken(oldRefreshToken)
	if err != nil {
		return nil, errors.ErrTokenInvalid.WithError(err)
	}

	// 验证refresh_token是否在Redis中存在
	exists, err := s.tokenStore.ValidateRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}
	if !exists {
		return nil, errors.ErrTokenInvalid.WithError(admin.ErrTokenInvalid) // Token已被撤销
	}

	// 删除旧的refresh_token
	_ = s.tokenStore.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID)
	_ = s.tokenStore.RemoveUserTokenID(ctx, claims.UserID, claims.TokenID)

	// 生成新的Token
	accessToken, _, err := jwt.Get().GenerateAccessToken(claims.UserID, claims.Username, "", claims.IsAdmin)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	refreshToken, refreshTokenID, err := jwt.Get().GenerateRefreshToken(claims.UserID, claims.Username, claims.IsAdmin)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	// 存储新的refresh_token到Redis
	if err := s.tokenStore.StoreRefreshToken(ctx, claims.UserID, refreshTokenID, refreshToken); err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}

	if err := s.tokenStore.StoreUserTokenID(ctx, claims.UserID, refreshTokenID); err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}