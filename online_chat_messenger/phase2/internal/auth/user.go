package auth

import (
	"errors"
	"online_chat_messenger/internal/chat"
)

// UserManager はユーザー管理のインターフェースです。
type UserManager interface {
	RegisterUser(token string, user chat.User) error
	FindUser(token string) (chat.User, error)
	DeleteUser(token string) error
}

// SimpleUserManager はUserManagerのシンプルな実装です。
type SimpleUserManager struct {
	users map[string]chat.User
}

// NewSimpleUserManager は新しいSimpleUserManagerを生成します。
func NewSimpleUserManager() *SimpleUserManager {
	return &SimpleUserManager{users: make(map[string]chat.User)}
}

// RegisterUser はユーザーを登録します。
func (m *SimpleUserManager) RegisterUser(token string, user chat.User) error {
	m.users[token] = user
	return nil
}

// FindUser は指定されたトークンのユーザーを返します。
func (m *SimpleUserManager) FindUser(token string) (chat.User, error) {
	user, ok := m.users[token]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// DeleteUser は指定されたトークンのユーザーを削除します。
func (m *SimpleUserManager) DeleteUser(token string) error {
	delete(m.users, token)
	return nil
}
