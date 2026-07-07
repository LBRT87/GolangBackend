package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/LBRT87/GolangBackend/services/auth-service/entity"
	"github.com/redis/go-redis/v9"
)

const (
	otpTTL          = 2 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type cacheRepository struct {
	client redis.UniversalClient
}

func NewCacheRepository(client redis.UniversalClient) entity.CacheRepository {
	return &cacheRepository{client: client}
}

func otpKey(email string) string {
	return fmt.Sprintf("otp:%s", email)
}

func refreshTokenKey(userID uint) string {
	return fmt.Sprintf("refresh_token:%d", userID)
}

func (r *cacheRepository) SetOTP(ctx context.Context, email string, code string) error {
	return r.client.Set(ctx, otpKey(email), code, otpTTL).Err()
}

func (r *cacheRepository) GetOTP(ctx context.Context, email string) (string, error) {
	code, err := r.client.Get(ctx, otpKey(email)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return code, err
}

func (r *cacheRepository) DeleteOTP(ctx context.Context, email string) error {
	return r.client.Del(ctx, otpKey(email)).Err()
}

func (r *cacheRepository) SetRefreshToken(ctx context.Context, userID uint, token string) error {
	return r.client.Set(ctx, refreshTokenKey(userID), token, refreshTokenTTL).Err()
}

func (r *cacheRepository) GetRefreshToken(ctx context.Context, userID uint) (string, error) {
	token, err := r.client.Get(ctx, refreshTokenKey(userID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return token, err
}

func (r *cacheRepository) DeleteRefreshToken(ctx context.Context, userID uint) error {
	return r.client.Del(ctx, refreshTokenKey(userID)).Err()
}
