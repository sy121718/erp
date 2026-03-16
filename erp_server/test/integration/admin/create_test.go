package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateAdmin 测试创建管理员
func TestCreateAdmin(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	createBody := map[string]string{
		"username": "testuser",
		"password": "test123456",
		"name":     "测试用户",
		"email":    "test@example.com",
		"phone":    "13800138000",
	}
	body, _ := json.Marshal(createBody)

	req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("创建管理员失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["code"].(float64) != 0 {
		t.Errorf("创建管理员返回错误: %v", resp)
		return
	}

	data := resp["data"].(map[string]interface{})
	testAdminID = int64(data["id"].(float64))
	t.Logf("创建管理员成功，ID: %d", testAdminID)
}

// TestCreateAdminDuplicate 测试创建重复用户名
func TestCreateAdminDuplicate(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("前置条件不满足，跳过测试")
	}

	createBody := map[string]string{
		"username": "testuser",
		"password": "test123456",
		"name":     "测试用户2",
	}
	body, _ := json.Marshal(createBody)

	req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["code"].(float64) == 0 {
		t.Error("重复用户名应该返回错误")
	}
}

// TestCreateAdminEmptyFields 测试空字段
func TestCreateAdminEmptyFields(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	testCases := []map[string]string{
		{"username": "", "password": "123456", "name": "test"},
		{"username": "test", "password": "", "name": "test"},
		{"username": "test", "password": "123456", "name": ""},
	}

	for i, tc := range testCases {
		body, _ := json.Marshal(tc)
		req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["code"].(float64) == 0 {
			t.Errorf("测试用例%d：空字段应该返回错误", i+1)
		}
	}
}

// TestCreateAdminSQLInjection 测试SQL注入
func TestCreateAdminSQLInjection(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	sqlPayloads := []string{
		"' OR '1'='1",
		"admin'; DROP TABLE admin; --",
	}

	for _, payload := range sqlPayloads {
		createBody := map[string]string{
			"username": payload,
			"password": "test123456",
			"name":     payload,
		}
		body, _ := json.Marshal(createBody)

		req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusInternalServerError {
			t.Errorf("SQL注入导致服务器错误: %s", payload)
		}
	}
}

// TestCreateAdminXSS 测试XSS
func TestCreateAdminXSS(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	xssPayloads := []string{
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert('xss')>",
	}

	for _, payload := range xssPayloads {
		createBody := map[string]string{
			"username": "xss_test",
			"password": "test123456",
			"name":     payload,
		}
		body, _ := json.Marshal(createBody)

		req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusInternalServerError {
			t.Errorf("XSS导致服务器错误: %s", payload)
		}
	}
}

// TestCreateAdminLongInput 测试超长输入
func TestCreateAdminLongInput(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	longString := strings.Repeat("a", 10000)

	createBody := map[string]string{
		"username": longString,
		"password": longString,
		"name":     longString,
	}
	body, _ := json.Marshal(createBody)

	req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusInternalServerError {
		t.Error("超长输入导致服务器崩溃")
	}
}

// TestCreateAdminWithoutToken 测试无token创建
func TestCreateAdminWithoutToken(t *testing.T) {
	createBody := map[string]string{
		"username": "noauth",
		"password": "test123456",
		"name":     "测试",
	}
	body, _ := json.Marshal(createBody)

	req := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token创建管理员应该失败")
	}
}
