package token

import (
	"context"
	"fmt"
	"time"

	"erp-server/pkg/config"
	"erp-server/pkg/redis"
)

const (
	// Redis key 前缀
	refreshTokenPrefix = "token:refresh:"
	accessTokenPrefix   = "token:access:"
	blacklistPrefix     = "token:blacklist:"
	userTokensPrefix    = "user:tokens:"
)

// TokenStore Token存储服务
type TokenStore struct {
	cfg *config.RedisConfig
}

var tokenStore *TokenStore

// Init 初始化Token存储
func Init(cfg *config.RedisConfig) {
	tokenStore = &TokenStore{cfg: cfg}
}

// Get 获取Token存储实例
func Get() *TokenStore {
	if tokenStore == nil {
		tokenStore = &TokenStore{cfg: config.DefaultRedisConfig()}
	}
	return tokenStore
}

// StoreRefreshToken 存储refresh_token
// Key: token:refresh:{user_id}:{token_id}
// Value: token_string
// TTL: 7天
func (s *TokenStore) StoreRefreshToken(ctx context.Context, userID int64, tokenID, token string) error {
	key := fmt.Sprintf("%s%d:%s", refreshTokenPrefix, userID, tokenID)
	return redis.Get().Set(ctx, key, token, s.cfg.GetRefreshTokenTTL()).Err()
}

// GetRefreshToken 获取refresh_token
func (s *TokenStore) GetRefreshToken(ctx context.Context, userID int64, tokenID string) (string, error) {
	key := fmt.Sprintf("%s%d:%s", refreshTokenPrefix, userID, tokenID)
	return redis.Get().Get(ctx, key).Result()
}

// DeleteRefreshToken 删除refresh_token（下线）
func (s *TokenStore) DeleteRefreshToken(ctx context.Context, userID int64, tokenID string) error {
	key := fmt.Sprintf("%s%d:%s", refreshTokenPrefix, userID, tokenID)
	return redis.Get().Del(ctx, key).Err()
}

// ValidateRefreshToken 验证refresh_token是否存在
func (s *TokenStore) ValidateRefreshToken(ctx context.Context, userID int64, tokenID string) (bool, error) {
	key := fmt.Sprintf("%s%d:%s", refreshTokenPrefix, userID, tokenID)
	exists, err := redis.Get().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// AddToBlacklist 将access_token加入黑名单
func (s *TokenStore) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := blacklistPrefix + token
	return redis.Get().Set(ctx, key, "1", expiration).Err()
}

// IsBlacklisted 检查token是否在黑名单中
func (s *TokenStore) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := blacklistPrefix + token
	exists, err := redis.Get().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// LogoutUser 用户下线（删除所有refresh_token）
// 通过删除用户的所有token实现强制下线
func (s *TokenStore) LogoutUser(ctx context.Context, userID int64) error {
	pattern := fmt.Sprintf("%s%d:*", refreshTokenPrefix, userID)
	iter := redis.Get().Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := redis.Get().Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// StoreUserTokenID 存储用户的token_id列表（用于管理多个设备登录）
func (s *TokenStore) StoreUserTokenID(ctx context.Context, userID int64, tokenID string) error {
	key := fmt.Sprintf("%s%d", userTokensPrefix, userID)
	// 使用SADD添加token_id到集合
	if err := redis.Get().SAdd(ctx, key, tokenID).Err(); err != nil {
		return err
	}
	// 设置过期时间
	return redis.Get().Expire(ctx, key, s.cfg.GetRefreshTokenTTL()).Err()
}

// GetUserTokenIDs 获取用户的所有token_id
func (s *TokenStore) GetUserTokenIDs(ctx context.Context, userID int64) ([]string, error) {
	key := fmt.Sprintf("%s%d", userTokensPrefix, userID)
	return redis.Get().SMembers(ctx, key).Result()
}

// RemoveUserTokenID 移除用户的某个token_id
func (s *TokenStore) RemoveUserTokenID(ctx context.Context, userID int64, tokenID string) error {
	key := fmt.Sprintf("%s%d", userTokensPrefix, userID)
	return redis.Get().SRem(ctx, key, tokenID).Err()
}