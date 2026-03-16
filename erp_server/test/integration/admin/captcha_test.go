package admin_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestGetCaptcha 测试获取验证码
func TestGetCaptcha(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/admin/captcha", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("获取验证码失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("解析响应失败: %v", err)
		return
	}

	if resp["code"].(float64) != 0 {
		t.Errorf("获取验证码返回错误: %v", resp)
	}

	data := resp["data"].(map[string]interface{})
	if data["captcha_id"] == nil || data["captcha_id"].(string) == "" {
		t.Error("验证码ID为空")
	}
}

// TestGetCaptchaMultipleRequests 测试多次请求验证码（压力测试）
func TestGetCaptchaMultipleRequests(t *testing.T) {
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/api/admin/captcha", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("第%d次获取验证码失败", i+1)
		}
	}
}

// TestCaptchaWithInvalidMethod 测试错误HTTP方法
func TestCaptchaWithInvalidMethod(t *testing.T) {
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/admin/captcha", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 应该返回404
		if w.Code == http.StatusOK {
			t.Errorf("验证码接口不应该支持%s方法", method)
		}
	}
}

// TestCaptchaWithLongPath 测试超长路径攻击
func TestCaptchaWithLongPath(t *testing.T) {
	// 使用有效的字符创建长路径
	longPath := "/api/admin/captcha/" + strings.Repeat("a", 1000)
	req := httptest.NewRequest("GET", longPath, nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 不应该崩溃，应该返回404
	if w.Code == http.StatusInternalServerError {
		t.Error("服务器不应该因超长路径而崩溃")
	}
}

// TestCaptchaWithSpecialChars 测试特殊字符
func TestCaptchaWithSpecialChars(t *testing.T) {
	// 使用 URL 编码的特殊字符
	specialPaths := []struct {
		path string
		desc string
	}{
		{"/api/admin/captcha?test=%3Cscript%3Ealert(1)%3C/script%3E", "XSS script tag"},
		{"/api/admin/captcha?test=%27OR%271%27%3D%271", "SQL injection"},
		{"/api/admin/captcha?test=..%2F..%2F..%2Fetc%2Fpasswd", "Path traversal"},
	}

	for _, tc := range specialPaths {
		req := httptest.NewRequest("GET", tc.path, nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 不应该崩溃，应该正常处理或返回错误
		if w.Code == http.StatusInternalServerError {
			t.Errorf("特殊字符路径导致服务器错误: %s", tc.desc)
		}
	}
}
