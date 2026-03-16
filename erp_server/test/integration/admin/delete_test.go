package admin_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDeleteAdmin 测试删除管理员
func TestDeleteAdmin(t *testing.T) {
	if adminToken == "" || testAdminID == 0 {
		t.Skip("没有token或管理员ID，跳过测试")
	}

	url := fmt.Sprintf("/api/admin/delete/%d", testAdminID)
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("删除管理员失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	t.Log("删除管理员成功")
	testAdminID = 0 // 已删除，重置ID
}

// TestDeleteAdminNonexistent 测试删除不存在的管理员
func TestDeleteAdminNonexistent(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("POST", "/api/admin/delete/99999999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["code"].(float64) == 0 {
		t.Error("删除不存在的管理员应该返回错误")
	}
}

// TestDeleteAdminInvalidID 测试无效ID
func TestDeleteAdminInvalidID(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	invalidIDs := []string{"abc", "-1", "0"}

	for _, id := range invalidIDs {
		url := fmt.Sprintf("/api/admin/delete/%s", id)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("无效ID '%s' 应该返回错误", id)
		}
	}
}

// TestDeleteAdminSelf 测试删除自己
func TestDeleteAdminSelf(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	// 获取自己的ID
	profileReq := httptest.NewRequest("GET", "/api/admin/profile", nil)
	profileReq.Header.Set("Authorization", "Bearer "+adminToken)
	profileW := httptest.NewRecorder()
	testRouter.ServeHTTP(profileW, profileReq)

	var profileResp map[string]interface{}
	json.Unmarshal(profileW.Body.Bytes(), &profileResp)

	if profileResp["code"].(float64) == 0 {
		data := profileResp["data"].(map[string]interface{})
		myID := int64(data["id"].(float64))

		url := fmt.Sprintf("/api/admin/delete/%d", myID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["code"].(float64) == 0 {
			t.Error("删除自己应该返回错误")
		}
	}
}

// TestDeleteAdminWithoutToken 测试无token删除
func TestDeleteAdminWithoutToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/admin/delete/1", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("无token删除管理员应该失败")
	}
}

// TestDeleteAdminTwice 测试重复删除
func TestDeleteAdminTwice(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	// 先创建一个临时管理员
	createBody := map[string]string{
		"username": "temp_delete_test",
		"password": "test123456",
		"name":     "临时删除测试",
	}
	body, _ := json.Marshal(createBody)

	createReq := httptest.NewRequest("POST", "/api/admin/create", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+adminToken)
	createW := httptest.NewRecorder()
	testRouter.ServeHTTP(createW, createReq)

	var createResp map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResp)

	if createResp["code"].(float64) == 0 {
		data := createResp["data"].(map[string]interface{})
		tempID := int64(data["id"].(float64))

		// 第一次删除
		url := fmt.Sprintf("/api/admin/delete/%d", tempID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		// 第二次删除
		req2 := httptest.NewRequest("POST", url, nil)
		req2.Header.Set("Authorization", "Bearer "+adminToken)
		w2 := httptest.NewRecorder()
		testRouter.ServeHTTP(w2, req2)

		var resp map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &resp)

		if resp["code"].(float64) == 0 {
			t.Error("重复删除应该返回错误")
		}
	}
}
