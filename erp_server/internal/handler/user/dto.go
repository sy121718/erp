package user

import (
	"erp-server/internal/domain/user"
)

// ========== 请求 DTO ==========

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Name   string `json:"name" binding:"omitempty"`
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone" binding:"omitempty"`
	Avatar string `json:"avatar" binding:"omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone" binding:"omitempty"`
}

// ListUserRequest 用户列表请求
type ListUserRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// ========== 响应 DTO ==========

// UserResponse 用户响应
type UserResponse struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Status        int    `json:"status"`
	Points        int64  `json:"points"`
	TotalPoints   int64  `json:"total_points"`
	UsedPoints    int64  `json:"used_points"`
	ExpireTime    string `json:"expire_time,omitempty"`
	LastLoginTime string `json:"last_login_time,omitempty"`
	CreateTime    string `json:"create_time"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
	User         *UserResponse `json:"user"`
}

// ListUserResponse 用户列表响应
type ListUserResponse struct {
	List     []*UserResponse `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// ========== 转换函数 ==========

// ToResponse Entity转响应DTO
func ToResponse(entity *user.User) *UserResponse {
	resp := &UserResponse{
		ID:          entity.ID,
		Username:    entity.Username,
		Name:        entity.Name,
		Avatar:      entity.Avatar,
		Email:       entity.Email,
		Phone:       entity.Phone,
		Status:      int(entity.Status),
		Points:      entity.Points,
		TotalPoints: entity.TotalPoints,
		UsedPoints:  entity.UsedPoints,
		CreateTime:  entity.CreateTime.Format("2006-01-02 15:04:05"),
	}

	if entity.ExpireTime != nil {
		resp.ExpireTime = entity.ExpireTime.Format("2006-01-02 15:04:05")
	}

	if entity.LastLoginTime != nil {
		resp.LastLoginTime = entity.LastLoginTime.Format("2006-01-02 15:04:05")
	}

	return resp
}

// ToListResponse 批量转换
func ToListResponse(entities []*user.User) []*UserResponse {
	list := make([]*UserResponse, len(entities))
	for i, e := range entities {
		list[i] = ToResponse(e)
	}
	return list
}
