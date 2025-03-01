package server

import (
	"net"
	"time"
)

type ClientManager struct {
	clients map[string]*net.UDPAddr // クライアントのリスト

	lastActive map[string]time.Time // クライアントごとの最終受信時刻

	timeout time.Duration // タイムアウト時刻
}

// コンストラクタ
func NewClientManager(timeout time.Duration) *ClientManager {
	cm := &ClientManager{
		clients:    make(map[string]*net.UDPAddr),
		lastActive: make(map[string]time.Time),
		timeout:    timeout,
	}

	return cm
}

func (cm *ClientManager) AddNewClient(addr *net.UDPAddr) {
	addrStr := addr.String()
	// 存在していなければ送信元を登録
	if _, exists := (*cm).clients[addrStr]; !exists {
		cm.clients[addrStr] = addr
		cm.lastActive[addrStr] = time.Now()
	}
}

func (cm *ClientManager) BroadCast(message string, senderAdder *net.UDPAddr, conn *net.UDPConn) {
	senderAdderStr := senderAdder.String()
	// clientsにmessageを送信する
	for _, addr := range cm.clients {
		if addr.String() == senderAdderStr {
			continue
		}
		// 送信処理もgoroutineにすることでループを効率的に回す
		go (*conn).WriteToUDP([]byte(message), addr) // 意識づけのためにあえて(*conn)としている
	}
}
