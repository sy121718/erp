package captcha

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Captcha 验证码信息
type Captcha struct {
	Code      string
	ExpiresAt time.Time
}

// Store 验证码存储接口
type Store interface {
	Set(id string, captcha *Captcha)
	Get(id string) *Captcha
	Delete(id string)
}

// MemoryStore 内存存储
type MemoryStore struct {
	mu       sync.RWMutex
	captchas map[string]*Captcha
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		captchas: make(map[string]*Captcha),
	}
	// 启动清理goroutine
	go store.cleanup()
	return store
}

// Set 设置验证码
func (s *MemoryStore) Set(id string, captcha *Captcha) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.captchas[id] = captcha
}

// Get 获取验证码
func (s *MemoryStore) Get(id string) *Captcha {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.captchas[id]
}

// Delete 删除验证码
func (s *MemoryStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.captchas, id)
}

// cleanup 定期清理过期验证码
func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, captcha := range s.captchas {
			if captcha.ExpiresAt.Before(now) {
				delete(s.captchas, id)
			}
		}
		s.mu.Unlock()
	}
}

// Config 验证码配置
type Config struct {
	Length     int           // 验证码长度
	ExpireTime time.Duration // 过期时间
	Width      int           // 图片宽度
	Height     int           // 图片高度
}

// CaptchaService 验证码服务
type CaptchaService struct {
	config *Config
	store  Store
}

var captchaService *CaptchaService

// Init 初始化验证码服务
func Init(cfg *Config) {
	if cfg == nil {
		cfg = &Config{
			Length:     6,
			ExpireTime: 1 * time.Minute,
			Width:      120,
			Height:     40,
		}
	}
	captchaService = &CaptchaService{
		config: cfg,
		store:  NewMemoryStore(),
	}
}

// Get 获取验证码服务
func Get() *CaptchaService {
	if captchaService == nil {
		Init(nil)
	}
	return captchaService
}

// Generate 生成验证码
func (s *CaptchaService) Generate() (id string, code string) {
	// 生成验证码ID
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id = fmt.Sprintf("%x", idBytes)

	// 生成数字验证码
	code = s.generateCode(s.config.Length)

	// 存储
	s.store.Set(id, &Captcha{
		Code:      code,
		ExpiresAt: time.Now().Add(s.config.ExpireTime),
	})

	return id, code
}

// GenerateWithPrefix 生成带前缀的验证码（区分业务场景）
func (s *CaptchaService) GenerateWithPrefix(prefix string) (id string, code string) {
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id = prefix + "_" + fmt.Sprintf("%x", idBytes)

	code = s.generateCode(s.config.Length)

	s.store.Set(id, &Captcha{
		Code:      code,
		ExpiresAt: time.Now().Add(s.config.ExpireTime),
	})

	return id, code
}

// Verify 验证验证码
func (s *CaptchaService) Verify(id, code string, clear bool) bool {
	captcha := s.store.Get(id)
	if captcha == nil {
		return false
	}

	// 检查是否过期
	if captcha.ExpiresAt.Before(time.Now()) {
		s.store.Delete(id)
		return false
	}

	// 验证码匹配
	if captcha.Code != code {
		return false
	}

	// 验证成功后删除
	if clear {
		s.store.Delete(id)
	}

	return true
}

// generateCode 生成指定长度的数字验证码
func (s *CaptchaService) generateCode(length int) string {
	digits := "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}

// GenerateAlphanumeric 生成字母数字混合验证码
func (s *CaptchaService) GenerateAlphanumeric() (id string, code string) {
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id = fmt.Sprintf("%x", idBytes)

	chars := "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789"
	codeBytes := make([]byte, s.config.Length)
	for i := 0; i < s.config.Length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		codeBytes[i] = chars[n.Int64()]
	}
	code = string(codeBytes)

	s.store.Set(id, &Captcha{
		Code:      string(code),
		ExpiresAt: time.Now().Add(s.config.ExpireTime),
	})

	return id, string(code)
}
