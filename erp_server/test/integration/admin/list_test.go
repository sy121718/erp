package admin_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListAdmins 测试获取管理员列表
func TestListAdmins(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("GET", "/api/admin/list?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("获取管理员列表失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		return
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["code"].(float64) != 0 {
		t.Errorf("获取管理员列表返回错误: %v", resp)
		return
	}

	data := resp["data"].(map[string]interface{})
	t.Logf("管理员总数: %v", data["total"])
}

// TestListAdminsWithKeyword 测试带关键字搜索
func TestListAdminsWithKeyword(t *testing.T) {
	if adminToken == "" {
		t.Skip("没有token，跳过测试")
	}

	req := httptest.NewRequest("GET", "/api/admin/list?page=1&page_size=10&keyword=sky", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("搜索管理员失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
	}
}

// TestListAdminsWithoutToken 测试无token
func TestListAdminsWithoutToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/admin/list?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	t.Logf("响应状态码: %d", w.Code)
	t.Logf("响应体: %s", w.Body.String())

	// 无token应该返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("无token应该返回401，实际返回: %d", w.Code)
	}
}