package admin

import (
	service "erp-server/internal/application/admin"
	"erp-server/internal/domain/admin"
	"erp-server/internal/handler/middleware"
	"erp-server/pkg/captcha"
	"erp-server/pkg/errors"
	"erp-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Handler 管理员处理器
type Handler struct {
	svc *service.AdminService
}

// NewHandler 创建处理器
func NewHandler(svc *service.AdminService) *Handler {
	return &Handler{svc: svc}
}

// Login 登录
// POST /api/admin/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	// 获取客户端信息
	ip := c.ClientIP()
	location := "" // 可以通过IP解析地理位置
	isp := ""

	resp, err := h.svc.Login(c.Request.Context(), &service.LoginRequest{
		Username:    req.Username,
		Password:    req.Password,
		CaptchaID:   req.CaptchaID,
		CaptchaCode: req.CaptchaCode,
		IP:          ip,
		Location:    location,
		ISP:         isp,
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, &LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    300, // 5分钟 = 300秒
		Admin:        ToResponse(resp.Admin),
	})
}

// Logout 登出
// POST /api/admin/logout
func (h *Handler) Logout(c *gin.Context) {
	adminID := middleware.GetUserID(c)

	// 从请求体获取refresh_token
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	c.ShouldBindJSON(&req) // 可选参数

	if err := h.svc.Logout(c.Request.Context(), adminID, req.RefreshToken); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// GetCaptcha 获取验证码
// GET /api/admin/captcha
func (h *Handler) GetCaptcha(c *gin.Context) {
	id, code := captcha.Get().Generate()
	utils.Success(c, &CaptchaResponse{
		CaptchaID: id,
		Code:      code, // 生产环境应删除
	})
}

// RefreshToken 刷新Token
// POST /api/admin/refresh-token
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
		"expires_in":    300, // 5分钟 = 300秒
	})
}

// GetProfile 获取当前管理员信息
// GET /api/admin/profile
func (h *Handler) GetProfile(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	entity, err := h.svc.GetByID(c.Request.Context(), adminID)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, ToResponse(entity))
}

// ChangePassword 修改自己的密码
// POST /api/admin/password/change
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	adminID := middleware.GetUserID(c)
	if err := h.svc.ChangePassword(c.Request.Context(), adminID, req.OldPassword, req.NewPassword); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// ResetPassword 重置管理员密码（仅超管）
// POST /api/admin/password/reset/:id
func (h *Handler) ResetPassword(c *gin.Context) {
	targetID, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	var req ResetPasswordRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := h.svc.ResetPassword(c.Request.Context(), operatorID, targetID, req.NewPassword); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// List 获取管理员列表
// GET /api/admin/list
func (h *Handler) List(c *gin.Context) {
	var req ListAdminRequest
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

	utils.Success(c, &ListAdminResponse{
		List:     ToListResponse(resp.List),
		Total:    resp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// Create 创建管理员（仅超管）
// POST /api/admin/create
func (h *Handler) Create(c *gin.Context) {
	// 检查是否是超管
	if !middleware.GetIsAdmin(c) {
		c.Error(errors.ErrNotAdmin)
		return
	}

	var req CreateAdminRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	operatorID := middleware.GetUserID(c)

	entity, err := h.svc.Create(c.Request.Context(), &service.CreateRequest{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		CreateBy: operatorID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, ToResponse(entity))
}

// Update 更新管理员（仅超管）
// POST /api/admin/update/:id
func (h *Handler) Update(c *gin.Context) {
	// 检查是否是超管
	if !middleware.GetIsAdmin(c) {
		c.Error(errors.ErrNotAdmin)
		return
	}

	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	var req UpdateAdminRequest
	if !middleware.BindJSON(c, &req) {
		return
	}

	operatorID := middleware.GetUserID(c)

	entity, err := h.svc.Update(c.Request.Context(), &service.UpdateRequest{
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

// Delete 删除管理员
// POST /api/admin/delete/:id
func (h *Handler) Delete(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	operatorID := middleware.GetUserID(c)

	if err := h.svc.Delete(c.Request.Context(), operatorID, id); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// Ban 禁用管理员
// POST /api/admin/ban/:id
func (h *Handler) Ban(c *gin.Context) {
	id, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	operatorID := middleware.GetUserID(c)

	if err := h.svc.Ban(c.Request.Context(), operatorID, id); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}

// Unban 解禁管理员
// POST /api/admin/unban/:id
func (h *Handler) Unban(c *gin.Context) {
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

// ForceLogout 强制下线（仅超管）
// POST /api/admin/force-logout/:id
func (h *Handler) ForceLogout(c *gin.Context) {
	targetID, ok := middleware.ParseID(c)
	if !ok {
		return
	}

	operatorID := middleware.GetUserID(c)

	// 检查是否是超管
	operator, err := h.svc.GetByID(c.Request.Context(), operatorID)
	if err != nil {
		c.Error(err)
		return
	}
	if !operator.IsAdmin {
		c.Error(errors.ErrNotAdmin)
		return
	}

	// 不能下线自己
	if operatorID == targetID {
		c.Error(admin.ErrCannotDisableSelf)
		return
	}

	if err := h.svc.ForceLogout(c.Request.Context(), targetID); err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, nil)
}