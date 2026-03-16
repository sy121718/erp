package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"erp-server/pkg/config"
)

// Signer 签名器
type Signer struct {
	secret     string
	expireTime int // 签名有效期(秒)
}

var signerInstance *Signer

// Init 初始化签名器
func Init(cfg *config.SignConfig) {
	if cfg == nil {
		cfg = config.DefaultSignConfig()
	}
	signerInstance = &Signer{
		secret:     cfg.GetSecret(),
		expireTime: cfg.GetExpireTime(),
	}
}

// Get 获取签名器实例
func Get() *Signer {
	if signerInstance == nil {
		signerInstance = &Signer{
			secret:     config.DefaultSignConfig().GetSecret(),
			expireTime: config.DefaultSignConfig().GetExpireTime(),
		}
	}
	return signerInstance
}

// SignParams 签名参数
type SignParams struct {
	Timestamp int64             `json:"timestamp"` // 时间戳(毫秒)
	Nonce     string            `json:"nonce"`     // 随机字符串
	Sign      string            `json:"sign"`      // 签名
	Params    map[string]string `json:"params"`    // 业务参数
}

// Generate 生成签名
// 签名规则:
// 1. 将所有参数按key升序排序
// 2. 拼接成 key1=value1&key2=value2 格式
// 3. 加上 timestamp 和 nonce
// 4. 使用 HMAC-SHA256 生成签名
func (s *Signer) Generate(params map[string]string, timestamp int64, nonce string) string {
	// 复制参数，避免修改原参数
	allParams := make(map[string]string)
	for k, v := range params {
		allParams[k] = v
	}
	allParams["timestamp"] = fmt.Sprintf("%d", timestamp)
	allParams["nonce"] = nonce

	// 获取所有key并排序
	keys := make([]string, 0, len(allParams))
	for k := range allParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(allParams[k])
	}
	signStr := builder.String()

	// HMAC-SHA256签名
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write([]byte(signStr))
	return hex.EncodeToString(h.Sum(nil))
}

// Verify 验证签名
func (s *Signer) Verify(params map[string]string, timestamp int64, nonce, sign string) bool {
	// 检查时间戳是否在有效期内
	now := time.Now().UnixMilli()
	if now-timestamp > int64(s.expireTime)*1000 {
		return false
	}

	// 检查时间戳是否在未来(允许1秒误差)
	if timestamp > now+1000 {
		return false
	}

	// 生成签名并比较
	expectedSign := s.Generate(params, timestamp, nonce)
	return hmac.Equal([]byte(sign), []byte(expectedSign))
}

// GenerateSignParams 生成带签名的完整参数
func (s *Signer) GenerateSignParams(params map[string]string) *SignParams {
	timestamp := time.Now().UnixMilli()
	nonce := s.generateNonce(16)
	sign := s.Generate(params, timestamp, nonce)

	return &SignParams{
		Timestamp: timestamp,
		Nonce:     nonce,
		Sign:      sign,
		Params:    params,
	}
}

// generateNonce 生成随机字符串
func (s *Signer) generateNonce(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[i%len(chars)]
	}
	return string(result)
}

// GenerateSimple 简单签名(不带参数)
func (s *Signer) GenerateSimple() (timestamp int64, nonce, sign string) {
	timestamp = time.Now().UnixMilli()
	nonce = s.generateNonce(16)
	sign = s.Generate(nil, timestamp, nonce)
	return
}

// VerifySimple 验证简单签名
func (s *Signer) VerifySimple(timestamp int64, nonce, sign string) bool {
	return s.Verify(nil, timestamp, nonce, sign)
}