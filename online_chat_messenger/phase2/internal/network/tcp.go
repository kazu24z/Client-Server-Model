package network

import (
	"fmt"
	"net"
)

// TCPServer はTCPサーバーを表します。
type TCPServer struct {
	listener net.Listener
}

// NewTCPServer は新しいTCPServerを生成します。
func NewTCPServer(port string) (*TCPServer, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("TCPサーバーの起動に失敗しました: %w", err)
	}
	return &TCPServer{listener: listener}, nil
}

// Start はTCPサーバーを起動し、クライアントからの接続を待ち受けます。
func (s *TCPServer) Start() error {
	fmt.Println("TCPサーバーを起動しました...")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("接続の受付に失敗しました: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection はクライアントとの接続を処理します。
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("クライアントが接続しました: %s\n", conn.RemoteAddr().String())

	// ここでクライアントとの通信処理を実装します。
}

// Close はTCPサーバーを停止します。
func (s *TCPServer) Close() error {
	return s.listener.Close()
}
