package chat

import "errors"

// RoomManager はチャットルーム管理のインターフェースです。
type RoomManager interface {
	CreateRoom(name, password string) (Room, error)
	FindRoom(name string) (Room, error)
	DeleteRoom(name string) error
}

// Room はチャットルームのインターフェースです。
type Room interface {
	GetName() string
	AddUser(user User, isHost bool) error
	RemoveUser(user User) error
	Broadcast(message string, sender User) error
	GetUsers() []User
}

// User はチャットルームのユーザーを表します。
type User interface {
	GetName() string
	GetToken() string
	GetAddress() string
	IsHost() bool
}

// SimpleRoomManager はRoomManagerのシンプルな実装です。
type SimpleRoomManager struct {
	rooms map[string]Room
}

// NewSimpleRoomManager は新しいSimpleRoomManagerを生成します。
func NewSimpleRoomManager() *SimpleRoomManager {
	return &SimpleRoomManager{rooms: make(map[string]Room)}
}

// CreateRoom は新しいチャットルームを作成します。
func (m *SimpleRoomManager) CreateRoom(name, password string) (Room, error) {
	if _, ok := m.rooms[name]; ok {
		return nil, errors.New("room already exists")
	}
	room := NewSimpleRoom(name, password)
	m.rooms[name] = room
	return room, nil
}

// FindRoom は指定された名前のチャットルームを返します。
func (m *SimpleRoomManager) FindRoom(name string) (Room, error) {
	room, ok := m.rooms[name]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}

// DeleteRoom は指定された名前のチャットルームを削除します。
func (m *SimpleRoomManager) DeleteRoom(name string) error {
	if _, ok := m.rooms[name]; !ok {
		return errors.New("room not found")
	}
	delete(m.rooms, name)
	return nil
}

// SimpleRoom はRoomのシンプルな実装です。
type SimpleRoom struct {
	name     string
	password string
	users    map[string]User
}

// NewSimpleRoom は新しいSimpleRoomを生成します。
func NewSimpleRoom(name, password string) *SimpleRoom {
	return &SimpleRoom{name: name, password: password, users: make(map[string]User)}
}

// GetName はチャットルームの名前を返します。
func (r *SimpleRoom) GetName() string {
	return r.name
}

// AddUser はチャットルームにユーザーを追加します。
func (r *SimpleRoom) AddUser(user User, isHost bool) error {
	r.users[user.GetToken()] = user
	return nil
}

// RemoveUser はチャットルームからユーザーを削除します。
func (r *SimpleRoom) RemoveUser(user User) error {
	delete(r.users, user.GetToken())
	return nil
}

// Broadcast はチャットルーム内の全ユーザーにメッセージを送信します。
func (r *SimpleRoom) Broadcast(message string, sender User) error {
	// ここでメッセージのブロードキャスト処理を実装します。
	return nil
}

// GetUsers はチャットルーム内の全ユーザーを返します。
func (r *SimpleRoom) GetUsers() []User {
	users := make([]User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

// SimpleUser はUserのシンプルな実装です。
type SimpleUser struct {
	name    string
	token   string
	address string
	isHost  bool
}

// NewUser は新しいSimpleUserを生成します。
func NewUser(name, token, address string) User {
	return &SimpleUser{name: name, token: token, address: address}
}

// GetName はユーザーの名前を返します。
func (u *SimpleUser) GetName() string {
	return u.name
}

// GetToken はユーザーのトークンを返します。
func (u *SimpleUser) GetToken() string {
	return u.token
}

// GetAddress はユーザーのアドレスを返します。
func (u *SimpleUser) GetAddress() string {
	return u.address
}

// IsHost はユーザーがホストかどうかを返します。
func (u *SimpleUser) IsHost() bool {
	return u.isHost
}
