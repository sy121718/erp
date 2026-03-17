package captcha

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// CaptchaType 验证码类型
type CaptchaType string

const (
	TypeDigit        CaptchaType = "digit"        // 纯数字
	TypeAlphanumeric CaptchaType = "alphanumeric"  // 字母数字混合
	TypeMath         CaptchaType = "math"          // 数学运算
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
	go store.cleanup()
	return store
}

func (s *MemoryStore) Set(id string, captcha *Captcha) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.captchas[id] = captcha
}

func (s *MemoryStore) Get(id string) *Captcha {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.captchas[id]
}

func (s *MemoryStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.captchas, id)
}

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

// GenerateByType 按类型生成验证码
// 返回 id, display（展示给用户的文本）, answer（正确答案）
func (s *CaptchaService) GenerateByType(captchaType CaptchaType) (id, display, answer string) {
	switch captchaType {
	case TypeAlphanumeric:
		id, code := s.GenerateAlphanumeric()
		return id, code, code
	case TypeMath:
		return s.GenerateMath()
	default:
		id, code := s.Generate()
		return id, code, code
	}
}

// Generate 生成纯数字验证码
func (s *CaptchaService) Generate() (id string, code string) {
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id = fmt.Sprintf("%x", idBytes)

	code = s.generateDigitCode(s.config.Length)

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

	code = s.generateDigitCode(s.config.Length)

	s.store.Set(id, &Captcha{
		Code:      code,
		ExpiresAt: time.Now().Add(s.config.ExpireTime),
	})

	return id, code
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
		Code:      code,
		ExpiresAt: time.Now().Add(s.config.ExpireTime),
	})

	return id, code
}

// GenerateMath 生成数学运算验证码
// 返回 id, expression（表达式，如 "3 + 5 = ?"）, answer（答案，如 "8"）
func (s *CaptchaService) GenerateMath() (id, expression, answer string) {
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	id = fmt.Sprintf("%x", idBytes)

	a := randInt(1, 20)
	b := randInt(1, 20)

	// 随机选择运算符
	ops := []string{"+", "-", "×"}
	opIdx := randInt(0, len(ops))
	op := ops[opIdx]

	var result int
	switch op {
	case "+":
		result = a + b
	case "-":
		if a < b {
			a, b = b, a
		}
		result = a - b
	case "×":
		a = randInt(1, 10)
		b = randInt(1, 10)
		result = a * b
	}

	expression = fmt.Sprintf("%d %s %d = ?", a, op, b)
	answer = fmt.Sprintf("%d", result)

	s.store.Set(id, &Captcha{
		Code:      answer,
		ExpiresAt: time.Now().Add(s.config.ExpireTime),
	})

	return id, expression, answer
}

// Verify 验证验证码
func (s *CaptchaService) Verify(id, code string, clear bool) bool {
	captcha := s.store.Get(id)
	if captcha == nil {
		return false
	}

	if captcha.ExpiresAt.Before(time.Now()) {
		s.store.Delete(id)
		return false
	}

	if captcha.Code != code {
		return false
	}

	if clear {
		s.store.Delete(id)
	}

	return true
}

func (s *CaptchaService) generateDigitCode(length int) string {
	digits := "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}

func randInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}
