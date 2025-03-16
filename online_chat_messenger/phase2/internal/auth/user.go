package auth

import (
	"errors"
	"fmt"
	"online_chat_messenger/internal/chat"
	"sync"
	"time"
)

// UserManager はユーザー管理のインターフェースです。
type UserManager interface {
	RegisterUser(token string, user chat.User) error
	FindUser(token string) (chat.User, error)
	DeleteUser(token string) error
}

// SimpleUserManager はUserManagerのシンプルな実装です。
type SimpleUserManager struct {
	users           map[string]chat.User
	lastActivityMap map[string]int64
	mutex           sync.RWMutex
	roomManager     chat.RoomManager
}

// NewSimpleUserManager は新しいSimpleUserManagerを生成します。
func NewSimpleUserManager() *SimpleUserManager {
	manager := &SimpleUserManager{
		users:           make(map[string]chat.User),
		lastActivityMap: make(map[string]int64),
		mutex:           sync.RWMutex{},
	}

	// 非アクティブユーザーを定期的に削除するゴルーチンを開始
	go manager.cleanupInactiveUsers()

	return manager
}

// RegisterUser はユーザーを登録します。
func (m *SimpleUserManager) RegisterUser(token string, user chat.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.users[token] = user
	m.lastActivityMap[token] = time.Now().Unix()
	return nil
}

// FindUser は指定されたトークンのユーザーを返します。
func (m *SimpleUserManager) FindUser(token string) (chat.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	user, ok := m.users[token]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// DeleteUser は指定されたトークンのユーザーを削除します。
func (m *SimpleUserManager) DeleteUser(token string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.users, token)
	delete(m.lastActivityMap, token)
	return nil
}

// UpdateActivity はユーザーの最終アクティビティ時間を更新します。
func (m *SimpleUserManager) UpdateActivity(token string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[token]; !exists {
		return errors.New("user not found")
	}

	m.lastActivityMap[token] = time.Now().Unix()
	return nil
}

// cleanupInactiveUsers は非アクティブなユーザーを定期的に削除します。
func (m *SimpleUserManager) cleanupInactiveUsers() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.removeInactiveUsers()
	}
}

// removeInactiveUsers は5分間アクティビティがなかったユーザーを削除します。
func (m *SimpleUserManager) removeInactiveUsers() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now().Unix()
	inactiveThreshold := int64(5 * 60) // 5分（秒単位）

	for token, lastActivity := range m.lastActivityMap {
		if now-lastActivity > inactiveThreshold {
			// ユーザーが所属するルームからも削除する必要がある
			if user, exists := m.users[token]; exists {
				// ルームマネージャーが設定されている場合、ルームからユーザーを削除
				if m.roomManager != nil {
					// すべてのルームをチェックし、ユーザーが所属していれば削除
					for _, room := range m.getRooms() {
						room.RemoveUser(user)
						fmt.Printf("非アクティブユーザー '%s' をルーム '%s' から削除しました\n",
							user.GetName(), room.GetName())
					}
				}
				fmt.Printf("非アクティブユーザー '%s' を削除します\n", user.GetName())
				delete(m.users, token)
				delete(m.lastActivityMap, token)
			}
		}
	}
}

// getRooms はルームマネージャーからすべてのルームを取得します
func (m *SimpleUserManager) getRooms() []chat.Room {
	if simpleRoomManager, ok := m.roomManager.(*chat.SimpleRoomManager); ok {
		return simpleRoomManager.GetAllRooms()
	}
	return []chat.Room{}
}

// SetRoomManager はルームマネージャーを設定します。
func (m *SimpleUserManager) SetRoomManager(roomManager chat.RoomManager) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.roomManager = roomManager
}
