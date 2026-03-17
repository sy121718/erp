package service

import (
	"context"
	"time"

	"erp-server/internal/domain/user"
	"erp-server/pkg/captcha"
	"erp-server/pkg/errors"
	"erp-server/pkg/jwt"
	"erp-server/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	repo       user.UserRepository
	tokenStore *token.TokenStore
}

// NewUserService 创建用户服务
func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{
		repo:       repo,
		tokenStore: token.Get(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string
	Password string
	Name     string
	Email    string
	Phone    string
	IP       string
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*user.User, error) {
	existing, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if existing != nil {
		return nil, errors.NewBadRequest(20001, user.ErrUsernameExists.Message())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	entity := &user.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Status:     user.StatusEnabled,
		Points:     0,
		TotalPoints: 0,
		UsedPoints:  0,
		UpdateTime: time.Now(),
	}

	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return entity, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username    string
	Password    string
	CaptchaID   string
	CaptchaCode string
	IP          string
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	User         *user.User
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 验证验证码
	if req.CaptchaID != "" && req.CaptchaCode != "" {
		if !captcha.Get().Verify(req.CaptchaID, req.CaptchaCode, true) {
			return nil, errors.NewBadRequest(400005, user.ErrCaptchaInvalid.Message())
		}
	}

	entity, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if entity == nil {
		return nil, errors.NewUnauthorized(401001, user.ErrPasswordIncorrect.Message())
	}

	if !entity.CanLogin() {
		return nil, errors.NewForbidden(403003, user.ErrAccountDisabled.Message())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password)); err != nil {
		return nil, errors.NewUnauthorized(401001, user.ErrPasswordIncorrect.Message())
	}

	entity.UpdateLastLogin(req.IP)
	s.repo.Update(ctx, entity)

	accessToken, _, err := jwt.Get().GenerateAccessToken(entity.ID, entity.Username, entity.Name, false, jwt.UserTypeUser)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	refreshToken, refreshTokenID, err := jwt.Get().GenerateRefreshToken(entity.ID, entity.Username, false, jwt.UserTypeUser)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	if err := s.tokenStore.StoreRefreshToken(ctx, entity.ID, refreshTokenID, refreshToken); err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}
	if err := s.tokenStore.StoreUserTokenID(ctx, entity.ID, refreshTokenID); err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         entity,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(ctx context.Context, userID int64, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	claims, err := jwt.Get().ParseRefreshToken(refreshToken)
	if err != nil {
		return nil
	}

	_ = s.tokenStore.DeleteRefreshToken(ctx, userID, claims.TokenID)
	_ = s.tokenStore.RemoveUserTokenID(ctx, userID, claims.TokenID)

	return nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(ctx context.Context, id int64) (*user.User, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	if entity == nil {
		return nil, errors.NewNotFound(404002, user.ErrUserNotFound.Message())
	}
	return entity, nil
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	ID       int64
	Name     string
	Email    string
	Phone    string
	Avatar   string
	UpdateBy int64
}

// Update 更新用户信息
func (s *UserService) Update(ctx context.Context, req *UpdateRequest) (*user.User, error) {
	entity, err := s.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	entity.Name = req.Name
	entity.Email = req.Email
	entity.Phone = req.Phone
	entity.Avatar = req.Avatar
	entity.UpdateBy = req.UpdateBy

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return entity, nil
}

// ChangePassword 修改自己的密码
func (s *UserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	entity, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(oldPassword)); err != nil {
		return errors.NewBadRequest(20004, user.ErrOldPasswordIncorrect.Message())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternal.WithError(err)
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
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
func (s *UserService) RefreshToken(ctx context.Context, oldRefreshToken string) (*RefreshTokenResponse, error) {
	claims, err := jwt.Get().ParseRefreshToken(oldRefreshToken)
	if err != nil {
		return nil, errors.ErrTokenInvalid.WithError(err)
	}

	if claims.UserType != jwt.UserTypeUser {
		return nil, errors.ErrTokenInvalid.WithError(err)
	}

	exists, err := s.tokenStore.ValidateRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		return nil, errors.ErrRedis.WithError(err)
	}
	if !exists {
		return nil, errors.ErrTokenInvalid
	}

	_ = s.tokenStore.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID)
	_ = s.tokenStore.RemoveUserTokenID(ctx, claims.UserID, claims.TokenID)

	accessToken, _, err := jwt.Get().GenerateAccessToken(claims.UserID, claims.Username, "", false, jwt.UserTypeUser)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

	refreshToken, refreshTokenID, err := jwt.Get().GenerateRefreshToken(claims.UserID, claims.Username, false, jwt.UserTypeUser)
	if err != nil {
		return nil, errors.ErrInternal.WithError(err)
	}

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

// ========== 管理员操作用户 ==========

// ListRequest 列表请求
type ListRequest struct {
	Page     int
	PageSize int
	Keyword  string
}

// ListResponse 列表响应
type ListResponse struct {
	List  []*user.User
	Total int64
}

// List 获取用户列表（管理员）
func (s *UserService) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	list, total, err := s.repo.List(ctx, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}
	return &ListResponse{
		List:  list,
		Total: total,
	}, nil
}

// AdminUpdateRequest 管理员更新用户请求
type AdminUpdateRequest struct {
	ID       int64
	Name     string
	Email    string
	Phone    string
	UpdateBy int64
}

// AdminUpdate 管理员更新用户信息
func (s *UserService) AdminUpdate(ctx context.Context, req *AdminUpdateRequest) (*user.User, error) {
	entity, err := s.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	entity.Name = req.Name
	entity.Email = req.Email
	entity.Phone = req.Phone
	entity.UpdateBy = req.UpdateBy

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, errors.ErrDatabase.WithError(err)
	}

	return entity, nil
}

// Ban 封禁用户（管理员）
func (s *UserService) Ban(ctx context.Context, targetID int64) error {
	_, err := s.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, targetID, user.StatusDisabled); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	_ = s.tokenStore.LogoutUser(ctx, targetID)

	return nil
}

// Unban 解封用户（管理员）
func (s *UserService) Unban(ctx context.Context, targetID int64) error {
	_, err := s.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, targetID, user.StatusEnabled); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	return nil
}

// ResetPassword 重置用户密码为 u123456（管理员）
func (s *UserService) ResetPassword(ctx context.Context, targetID int64) error {
	_, err := s.GetByID(ctx, targetID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("u123456"), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternal.WithError(err)
	}

	if err := s.repo.UpdatePassword(ctx, targetID, string(hashedPassword)); err != nil {
		return errors.ErrDatabase.WithError(err)
	}

	_ = s.tokenStore.LogoutUser(ctx, targetID)

	return nil
}
