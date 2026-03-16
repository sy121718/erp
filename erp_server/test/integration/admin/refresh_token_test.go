package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRefreshToken 测试刷新Token
func TestRefreshToken(t *testing.T) {
	// 先登录获取token
	captchaID, captchaCode := getCaptcha(t)

	loginBody := map[string]string{
		"username":   "sky",
		"password":   "123456",
		"captcha_id": captchaID,
		"captcha":    captchaCode,
	}
	body, _ := json.Marshal(loginBody)

	req := httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Skip("登录失败，跳过刷新Token测试")
		return
	}

	var loginResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Skip("解析登录响应失败")
		return
	}

	if loginResp["code"] == nil || loginResp["code"].(float64) != 0 {
		t.Skip("登录返回错误")
		return
	}

	data := loginResp["data"].(map[string]interface{})
	testRefreshToken := data["refresh_token"].(string)

	if testRefreshToken == "" {
		t.Skip("登录未返回refresh_token，跳过测试")
		return
	}

	refreshBody := map[string]string{
		"refresh_token": testRefreshToken,
	}
	body, _ = json.Marshal(refreshBody)

	req = httptest.NewRequest("POST", "/api/admin/refresh-token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("刷新Token失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("解析响应失败: %v", err)
		return
	}

	if resp["code"] == nil || resp["code"].(float64) != 0 {
		t.Errorf("刷新Token返回错误: %v", resp)
		return
	}

	respData := resp["data"].(map[string]interface{})
	newAccessToken := respData["access_token"].(string)
	newRefreshToken := respData["refresh_token"].(string)

	if newAccessToken == "" {
		t.Error("刷新Token后未返回access_token")
	}
	if newRefreshToken == "" {
		t.Error("刷新Token后未返回refresh_token")
	}

	t.Logf("刷新Token成功，新access_token: %s...", newAccessToken[:20])
}

// TestRefreshTokenInvalid 测试无效Token刷新
func TestRefreshTokenInvalid(t *testing.T) {
	invalidTokens := []string{
		"invalid_token",
		"",
	}

	for _, token := range invalidTokens {
		refreshBody := map[string]string{"refresh_token": token}
		body, _ := json.Marshal(refreshBody)

		req := httptest.NewRequest("POST", "/api/admin/refresh-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 无效token应该返回错误
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
			t.Errorf("无效token '%s' 刷新应该失败", token)
		}
	}
}

// TestRefreshTokenEmptyBody 测试空请求体
func TestRefreshTokenEmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/admin/refresh-token", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 空请求体应该返回错误
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
		t.Error("空请求体应该返回错误")
	}
}

// TestRefreshTokenInvalidJSON 测试无效JSON
func TestRefreshTokenInvalidJSON(t *testing.T) {
	invalidJSONs := [][]byte{
		[]byte("{invalid}"),
		[]byte("null"),
	}

	for _, invalidJSON := range invalidJSONs {
		req := httptest.NewRequest("POST", "/api/admin/refresh-token", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 测试目的是确保服务器不会崩溃
		t.Logf("无效JSON测试状态码: %d", w.Code)
	}
}

// TestRefreshTokenMissingField 测试缺少字段
func TestRefreshTokenMissingField(t *testing.T) {
	bodies := []map[string]interface{}{
		{},
		{"wrong_field": "value"},
	}

	for _, body := range bodies {
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/admin/refresh-token", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if w.Code == http.StatusOK && resp["code"] != nil && resp["code"].(float64) == 0 {
			t.Error("缺少refresh_token字段应该返回错误")
		}
	}
}