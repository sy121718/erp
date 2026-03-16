package admin_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestResetPassword 测试重置密码
func TestResetPassword(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	resetBody := map[string]string{
		"new_password": "newtest123456",
	}
	body, _ := json.Marshal(resetBody)

	url := fmt.Sprintf("/api/admin/password/reset/%d", testAdminID)
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Logf("重置密码状态码: %d", w.Code)
}

// TestChangePassword 测试修改密码
func TestChangePassword(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	changeBody := map[string]string{
		"old_password": "123456",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(changeBody)

	req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Logf("修改密码状态码: %d", w.Code)
}

// TestChangePasswordWrongOld 测试错误原密码
func TestChangePasswordWrongOld(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	changeBody := map[string]string{
		"old_password": "wrongoldpassword",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(changeBody)

	req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 应该返回错误
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("错误原密码应该返回错误")
	}
}

// TestChangePasswordEmptyFields 测试空字段
func TestChangePasswordEmptyFields(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	testCases := []map[string]string{
		{"old_password": "", "new_password": "123456"},
		{"old_password": "123456", "new_password": ""},
	}

	for i, tc := range testCases {
		body, _ := json.Marshal(tc)
		req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
			t.Errorf("测试用例%d：空字段应该返回错误", i+1)
		}
	}
}

// TestChangePasswordTooShort 测试密码太短
func TestChangePasswordTooShort(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	changeBody := map[string]string{
		"old_password": "123456",
		"new_password": "123",
	}
	body, _ := json.Marshal(changeBody)

	req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("密码太短应该返回错误")
	}
}

// TestPasswordSQLInjection 测试SQL注入
func TestPasswordSQLInjection(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	sqlPayloads := []string{
		"' OR '1'='1",
		"123'; DROP TABLE admin; --",
	}

	for _, payload := range sqlPayloads {
		changeBody := map[string]string{
			"old_password": payload,
			"new_password": payload,
		}
		body, _ := json.Marshal(changeBody)

		req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 确保不会崩溃
		t.Logf("SQL注入测试状态码: %d", w.Code)
	}
}

// TestPasswordLongInput 测试超长密码
func TestPasswordLongInput(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	longPassword := strings.Repeat("a", 10000)

	changeBody := map[string]string{
		"old_password": "123456",
		"new_password": longPassword,
	}
	body, _ := json.Marshal(changeBody)

	req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Logf("超长密码测试状态码: %d", w.Code)
}

// TestPasswordWithoutToken 测试无token
func TestPasswordWithoutToken(t *testing.T) {
	changeBody := map[string]string{
		"old_password": "123456",
		"new_password": "newpassword",
	}
	body, _ := json.Marshal(changeBody)

	req := httptest.NewRequest("POST", "/api/admin/password/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token修改密码应该失败")
	}
}