package auth

import (
	"context"
	"errors"
)

// ErrInvalidToken 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

// TokenInfo 令牌信息
type TokenInfo interface {
	// GetAccessToken 获取访问令牌
	GetAccessToken() string
	// GetTokenType 获取令牌类型
	GetTokenType() string
	// GetExpiresAt 获取令牌到期时间戳
	GetExpiresAt() int64
	// EncodeToJSON JSON编码
	EncodeToJSON() ([]byte, error)
}

// Auther 认证接口
type Auther interface {
	// GenerateToken 生成令牌
	GenerateToken(ctx context.Context, userID string) (TokenInfo, error)

	// DestroyToken 销毁令牌
	DestroyToken(ctx context.Context, accessToken string) error

	// ParseUserID 解析用户ID
	ParseUserID(ctx context.Context, accessToken string) (string, error)

	// Release 释放资源
	Release() error
}
