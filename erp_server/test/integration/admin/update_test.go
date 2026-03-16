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

// TestUpdateAdmin 测试更新管理员
func TestUpdateAdmin(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	updateBody := map[string]string{
		"name":  "测试用户-已更新",
		"email": "updated@example.com",
		"phone": "13900139000",
	}
	body, _ := json.Marshal(updateBody)

	url := fmt.Sprintf("/api/admin/update/%d", testAdminID)
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("更新管理员失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["code"].(float64) != 0 {
		t.Errorf("更新管理员返回错误: %v", resp)
		return
	}

	t.Logf("更新管理员成功")
}

// TestUpdateAdminNonexistent 测试更新不存在的管理员
func TestUpdateAdminNonexistent(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	updateBody := map[string]string{
		"name":  "不存在",
		"email": "none@example.com",
	}
	body, _ := json.Marshal(updateBody)

	req := httptest.NewRequest("POST", "/api/admin/update/99999999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 测试目的是确保服务器不会崩溃
	t.Logf("更新不存在管理员状态码: %d", w.Code)
}

// TestUpdateAdminInvalidID 测试无效ID
func TestUpdateAdminInvalidID(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	invalidIDs := []string{"abc", "-1", "0"}

	for _, id := range invalidIDs {
		url := fmt.Sprintf("/api/admin/update/%s", id)
		updateBody := map[string]string{"name": "test"}
		body, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("无效ID '%s' 应该返回错误", id)
		}
	}
}

// TestUpdateAdminSQLInjection 测试SQL注入
func TestUpdateAdminSQLInjection(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	sqlPayloads := []string{
		"' OR '1'='1",
		"test'; DROP TABLE admin; --",
	}

	for _, payload := range sqlPayloads {
		updateBody := map[string]string{
			"name":  payload,
			"email": payload,
		}
		body, _ := json.Marshal(updateBody)

		url := fmt.Sprintf("/api/admin/update/%d", testAdminID)
		req := httptest.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusInternalServerError {
			t.Errorf("SQL注入导致服务器错误: %s", payload)
		}
	}
}

// TestUpdateAdminLongInput 测试超长输入
func TestUpdateAdminLongInput(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	longString := strings.Repeat("a", 10000)

	updateBody := map[string]string{
		"name":  longString,
		"email": "test@example.com",
	}
	body, _ := json.Marshal(updateBody)

	url := fmt.Sprintf("/api/admin/update/%d", testAdminID)
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusInternalServerError {
		t.Error("超长输入导致服务器崩溃")
	}
}

// TestUpdateAdminWithoutToken 测试无token更新
func TestUpdateAdminWithoutToken(t *testing.T) {
	updateBody := map[string]string{"name": "test"}
	body, _ := json.Marshal(updateBody)

	req := httptest.NewRequest("POST", "/api/admin/update/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token更新管理员应该失败")
	}
}
