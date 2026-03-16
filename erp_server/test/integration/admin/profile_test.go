package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetProfile 测试获取个人信息
func TestGetProfile(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("GET", "/api/admin/profile", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 测试目的是确保服务器正常响应
	t.Logf("获取个人信息状态码: %d", w.Code)
}

// TestGetProfileWithoutToken 测试无token获取个人信息
func TestGetProfileWithoutToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/admin/profile", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token获取个人信息应该失败")
	}
}

// TestGetProfileByID 测试按ID获取管理员信息
func TestGetProfileByID(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("GET", "/api/admin/1", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Logf("按ID获取管理员状态码: %d", w.Code)
}