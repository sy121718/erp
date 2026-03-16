package admin

import (
	"erp-server/internal/domain/admin"
)

// ========== 请求 DTO ==========

// LoginRequest 登录请求
type LoginRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
}

// UpdateAdminRequest 更新管理员请求
type UpdateAdminRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone" binding:"omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// ResetPasswordRequest 重置密码请求（超管使用）
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// ListAdminRequest 管理员列表请求
type ListAdminRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// ========== 响应 DTO ==========

// AdminResponse 管理员响应
type AdminResponse struct {
	ID                int64  `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Avatar            string `json:"avatar"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Status            int    `json:"status"`
	IsAdmin           bool   `json:"is_admin"`
	LoginFailureCount int    `json:"login_failure_count"`
	LastLoginTime     string `json:"last_login_time,omitempty"`
	CreateTime        string `json:"create_time"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	ExpiresIn    int             `json:"expires_in"` // access_token过期时间(秒)
	Admin        *AdminResponse `json:"admin"`
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	Code      string `json:"code,omitempty"` // 仅开发环境返回
}

// ListAdminResponse 管理员列表响应
type ListAdminResponse struct {
	List     []*AdminResponse `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// ========== 转换函数 ==========

// ToResponse Entity转响应DTO
func ToResponse(entity *admin.Admin) *AdminResponse {
	resp := &AdminResponse{
		ID:                entity.ID,
		Username:          entity.Username,
		Name:              entity.Name,
		Avatar:            entity.Avatar,
		Email:             entity.Email,
		Phone:             entity.Phone,
		Status:            int(entity.Status),
		IsAdmin:           entity.IsAdmin,
		LoginFailureCount: entity.LoginFailureCount,
		CreateTime:        entity.CreateTime.Format("2006-01-02 15:04:05"),
	}

	if entity.LastLoginTime != nil {
		resp.LastLoginTime = entity.LastLoginTime.Format("2006-01-02 15:04:05")
	}

	return resp
}

// ToListResponse 批量转换
func ToListResponse(entities []*admin.Admin) []*AdminResponse {
	list := make([]*AdminResponse, len(entities))
	for i, e := range entities {
		list[i] = ToResponse(e)
	}
	return list
}