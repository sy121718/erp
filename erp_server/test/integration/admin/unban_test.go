package admin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUnbanAdmin 测试解禁管理员
func TestUnbanAdmin(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	url := fmt.Sprintf("/api/admin/unban/%d", testAdminID)
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Logf("解禁管理员状态码: %d", w.Code)
}

// TestUnbanAdminNonexistent 测试解禁不存在的管理员
func TestUnbanAdminNonexistent(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("POST", "/api/admin/unban/99999999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 测试目的是确保服务器不会崩溃
	t.Logf("解禁不存在管理员状态码: %d", w.Code)
}

// TestUnbanAdminInvalidID 测试无效ID
func TestUnbanAdminInvalidID(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	invalidIDs := []string{"abc", "-1", "0"}

	for _, id := range invalidIDs {
		url := fmt.Sprintf("/api/admin/unban/%s", id)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("无效ID '%s' 应该返回错误", id)
		}
	}
}

// TestUnbanAdminWithoutToken 测试无token解禁
func TestUnbanAdminWithoutToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/admin/unban/1", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token解禁管理员应该失败")
	}
}

// TestBanUnbanFlow 测试禁用解禁流程
func TestBanUnbanFlow(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	// 禁用
	banURL := fmt.Sprintf("/api/admin/ban/%d", testAdminID)
	banReq := httptest.NewRequest("POST", banURL, nil)
	banReq.Header.Set("Authorization", "Bearer "+adminToken)
	banW := httptest.NewRecorder()
	testRouter.ServeHTTP(banW, banReq)

	t.Logf("禁用状态码: %d", banW.Code)

	// 解禁
	unbanURL := fmt.Sprintf("/api/admin/unban/%d", testAdminID)
	unbanReq := httptest.NewRequest("POST", unbanURL, nil)
	unbanReq.Header.Set("Authorization", "Bearer "+adminToken)
	unbanW := httptest.NewRecorder()
	testRouter.ServeHTTP(unbanW, unbanReq)

	t.Logf("解禁状态码: %d", unbanW.Code)
}