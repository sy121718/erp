package admin_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBanAdmin 测试禁用管理员
func TestBanAdmin(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	url := fmt.Sprintf("/api/admin/ban/%d", testAdminID)
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("禁用管理员失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	t.Log("禁用管理员成功")
}

// TestBanAdminNonexistent 测试禁用不存在的管理员
func TestBanAdminNonexistent(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("POST", "/api/admin/ban/99999999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["code"].(float64) == 0 {
		t.Error("禁用不存在的管理员应该返回错误")
	}
}

// TestBanAdminInvalidID 测试无效ID
func TestBanAdminInvalidID(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	invalidIDs := []string{"abc", "-1", "0"}

	for _, id := range invalidIDs {
		url := fmt.Sprintf("/api/admin/ban/%s", id)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("无效ID '%s' 应该返回错误", id)
		}
	}
}

// TestBanAdminSelf 测试禁用自己
func TestBanAdminSelf(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	// 先获取自己的ID
	profileReq := httptest.NewRequest("GET", "/api/admin/profile", nil)
	profileReq.Header.Set("Authorization", "Bearer "+adminToken)
	profileW := httptest.NewRecorder()
	testRouter.ServeHTTP(profileW, profileReq)

	var profileResp map[string]interface{}
	json.Unmarshal(profileW.Body.Bytes(), &profileResp)

	if profileResp["code"].(float64) == 0 {
		data := profileResp["data"].(map[string]interface{})
		myID := int64(data["id"].(float64))

		url := fmt.Sprintf("/api/admin/ban/%d", myID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["code"].(float64) == 0 {
			t.Error("禁用自己应该返回错误")
		}
	}
}

// TestBanAdminWithoutToken 测试无token禁用
func TestBanAdminWithoutToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/admin/ban/1", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token禁用管理员应该失败")
	}
}
