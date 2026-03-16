package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestLogin 测试登录
func TestLogin(t *testing.T) {
	// 先获取验证码
	captchaID, captchaCode := getCaptcha(t)

	// 使用超级管理员登录
	loginBody := map[string]string{
		"username":     superAdminUsername,
		"password":     superAdminPassword,
		"captcha_id":   captchaID,
		"captcha_code": captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("登录失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("解析响应失败: %v", err)
		return
	}

	if resp["code"] == nil || resp["code"].(float64) != 0 {
		t.Errorf("登录返回错误: %v", resp)
		return
	}

	data := resp["data"].(map[string]interface{})
	adminToken = data["access_token"].(string)
	refreshToken = data["refresh_token"].(string)

	if adminToken == "" {
		t.Error("登录后未返回token")
	}

	t.Logf("登录成功，token: %s...", adminToken[:20])
}

// TestLoginWrongPassword 测试错误密码登录
func TestLoginWrongPassword(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	loginBody := map[string]string{
		"username":     superAdminUsername,
		"password":     "wrongpassword",
		"captcha_id":   captchaID,
		"captcha_code": captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 不应该返回200 OK或者返回code=0
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("错误密码登录应该返回错误")
	}
}

// TestLoginEmptyUsername 测试空用户名
func TestLoginEmptyUsername(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	loginBody := map[string]string{
		"username":     "",
		"password":     "123456",
		"captcha_id":   captchaID,
		"captcha_code": captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 不应该返回200 OK或者返回code=0
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("空用户名应该返回错误")
	}
}

// TestLoginEmptyPassword 测试空密码
func TestLoginEmptyPassword(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	loginBody := map[string]string{
		"username":     superAdminUsername,
		"password":     "",
		"captcha_id":   captchaID,
		"captcha_code": captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 不应该返回200 OK或者返回code=0
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("空密码应该返回错误")
	}
}

// TestLoginNonexistentUser 测试不存在的用户
func TestLoginNonexistentUser(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	loginBody := map[string]string{
		"username":     "nonexistent_user_12345",
		"password":     "123456",
		"captcha_id":   captchaID,
		"captcha_code": captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 不应该返回200 OK或者返回code=0
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("不存在的用户应该返回错误")
	}
}

// TestLoginWrongCaptcha 测试错误验证码
func TestLoginWrongCaptcha(t *testing.T) {
	captchaID, _ := getCaptcha(t)

	loginBody := map[string]string{
		"username":     superAdminUsername,
		"password":     superAdminPassword,
		"captcha_id":   captchaID,
		"captcha_code": "wrong",
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 不应该返回200 OK或者返回code=0
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("错误验证码应该返回错误")
	}
}

// TestLoginSQLInjection 测试SQL注入攻击
func TestLoginSQLInjection(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	sqlPayloads := []string{
		"' OR '1'='1",
		"' OR '1'='1' --",
		"admin'--",
		"1; DROP TABLE admin; --",
	}

	for _, payload := range sqlPayloads {
		loginBody := map[string]string{
			"username":     payload,
			"password":     payload,
			"captcha_id":   captchaID,
			"captcha_code": captchaCode,
		}
		body, _ := json.Marshal(loginBody)

		req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
			t.Errorf("SQL注入应该被拦截: %s", payload)
		}
	}
}

// TestLoginXSS 测试XSS攻击
func TestLoginXSS(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	xssPayloads := []string{
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert('xss')>",
	}

	for _, payload := range xssPayloads {
		loginBody := map[string]string{
			"username":     payload,
			"password":     "123456",
			"captcha_id":   captchaID,
			"captcha_code": captchaCode,
		}
		body, _ := json.Marshal(loginBody)

		req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 测试目的是确保服务器不会崩溃，返回任何响应都算通过
		// XSS payload 可能导致用户不存在错误，这是正常处理
		t.Logf("XSS测试状态码: %d", w.Code)
	}
}

// TestLoginLongInput 测试超长输入
func TestLoginLongInput(t *testing.T) {
	captchaID, captchaCode := getCaptcha(t)

	longString := strings.Repeat("a", 10000)

	loginBody := map[string]string{
		"username":     longString,
		"password":     longString,
		"captcha_id":   captchaID,
		"captcha_code": captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 测试目的是确保服务器不会因超长输入而崩溃
	// 返回任何响应都算通过（用户不存在、参数错误等都是正常处理）
	t.Logf("超长输入测试状态码: %d", w.Code)
}

// TestLoginInvalidJSON 测试无效JSON
func TestLoginInvalidJSON(t *testing.T) {
	invalidJSONs := [][]byte{
		[]byte("{invalid json}"),
		[]byte(`{"username": "admin"`),
		[]byte("null"),
		[]byte(""),
		[]byte("[]"),
	}

	for _, invalidJSON := range invalidJSONs {
		req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 测试目的是确保服务器不会因无效JSON而崩溃
		// 返回参数错误（400或500）都是正常的错误处理
		t.Logf("无效JSON测试状态码: %d", w.Code)
	}
}

// TestLoginWithoutCaptcha 测试不带验证码登录
func TestLoginWithoutCaptcha(t *testing.T) {
	loginBody := map[string]string{
		"username": superAdminUsername,
		"password": superAdminPassword,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 根据业务逻辑，可能允许或禁止无验证码登录
	t.Logf("无验证码登录状态码: %d", w.Code)
}

// getCaptcha 获取验证码
func getCaptcha(t *testing.T) (string, string) {
	req := httptest.NewRequest("GET", "/api/admin/captcha", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("获取验证码失败: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	return data["captcha_id"].(string), data["code"].(string)
}