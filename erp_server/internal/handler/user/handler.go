package user

import (
	service "erp-server/internal/application/user"
	"erp-server/internal/handler/middleware"
	"erp-server/pkg/captcha"
	"erp-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Handler 用户处理器
type Handler struct {
	svc *service.UserService
}

// NewHandler 创建处理器
func NewHandler(svc *service.UserService) *Handler {
	return &Handler{svc: svc}
}

// GetCaptcha 获取验证码
// GET /api/user/captcha
func (h *Handler) GetCaptcha(c *gin.Context) {
	id, code := captcha.Get().Generate()
	utils.Success(c, gin.H{
		"captcha_id": id,
		"code":       code,
	})
}

// Register 用户注册
// POST /api/user/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	entity, err := h.svc.Register(c.Request.Context(), &service.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		IP:       c.ClientIP(),
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, ToResponse(entity))
}

// Login 用户登录
// POST /api/user/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), &service.LoginRequest{
		Username:    req.Username,
		Password:    req.Password,
		CaptchaID:   req.CaptchaID,
		CaptchaCode: req.CaptchaCode,
		IP:          c.ClientIP(),
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, &LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    300,
		User:         ToResponse(resp.User),
	})
}

// Logout 用户登出
// POST /api/user/logout
func (h *Handler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.Logout(c.Request.Context(), userID, req.RefreshToken); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// GetProfile 获取当前用户信息
// GET /api/user/profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	entity, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, ToResponse(entity))
}

// UpdateProfile 更新个人信息
// POST /api/user/update
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	entity, err := h.svc.Update(c.Request.Context(), &service.UpdateRequest{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		UpdateBy: userID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, ToResponse(entity))
}

// ChangePassword 修改密码
// POST /api/user/password/change
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.svc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// RefreshToken 刷新Token
// POST /api/user/refresh-token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	resp, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    300,
	})
}

// ========== 管理员操作用户接口 ==========

// AdminListUsers 管理员查询用户列表
// GET /api/admin/user/list
func (h *Handler) AdminListUsers(c *gin.Context) {
	var req ListUserRequest
	if !middleware.BindQuery(c, &req) {
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	resp, err := h.svc.List(c.Request.Context(), &service.ListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, &ListUserResponse{
		List:     ToListResponse(resp.List),
		Total:    resp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// AdminGetUser 管理员查询用户详情
// GET /api/admin/user/:id
func (h *Handler) AdminGetUser(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	entity, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, ToResponse(entity))
}

// AdminUpdateUser 管理员更新用户信息
// POST /api/admin/user/update/:id
func (h *Handler) AdminUpdateUser(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	var req AdminUpdateUserRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	operatorID := middleware.GetUserID(c)
	entity, err := h.svc.AdminUpdate(c.Request.Context(), &service.AdminUpdateRequest{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		UpdateBy: operatorID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, ToResponse(entity))
}

// AdminBanUser 管理员封禁用户
// POST /api/admin/user/ban/:id
func (h *Handler) AdminBanUser(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	if err := h.svc.Ban(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// AdminUnbanUser 管理员解封用户
// POST /api/admin/user/unban/:id
func (h *Handler) AdminUnbanUser(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	if err := h.svc.Unban(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// AdminResetPassword 管理员重置用户密码为 u123456
// POST /api/admin/user/password/reset/:id
func (h *Handler) AdminResetPassword(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	if err := h.svc.ResetPassword(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}
