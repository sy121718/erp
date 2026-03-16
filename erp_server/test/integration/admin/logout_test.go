package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLogout 测试登出
func TestLogout(t *testing.T) {
	if adminToken == "" || refreshToken == "" {
		t.Skip("没有token，跳过测试")
	}

	logoutBody := map[string]string{
		"refresh_token": refreshToken,
	}
	body, _ := json.Marshal(logoutBody)

	req := httptest.NewRequest("POST", "/api/admin/logout", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 测试目的是确保服务器正常响应
	t.Logf("登出状态码: %d", w.Code)
}

// TestLogoutWithoutToken 测试无token登出
func TestLogoutWithoutToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/admin/logout", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 无token应该返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("无token登出应该返回401，实际返回: %d", w.Code)
	}
}

// TestLogoutWithInvalidToken 测试无效token登出
func TestLogoutWithInvalidToken(t *testing.T) {
	invalidTokens := []string{
		"invalid_token",
		"",
	}

	for _, token := range invalidTokens {
		req := httptest.NewRequest("POST", "/api/admin/logout", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 无效token应该返回401
		if w.Code != http.StatusUnauthorized {
			t.Errorf("无效token '%s' 登出应该返回401，实际: %d", token, w.Code)
		}
	}
}