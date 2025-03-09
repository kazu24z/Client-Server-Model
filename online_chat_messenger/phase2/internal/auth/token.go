package auth

import "github.com/google/uuid"

// GenerateToken は新しいトークンを生成します。
func GenerateToken() string {
	return uuid.New().String()
}
