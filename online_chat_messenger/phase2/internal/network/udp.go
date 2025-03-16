package network

import (
	"fmt"
	"net"
	"online_chat_messenger/internal/auth"
	"online_chat_messenger/internal/chat"
	"online_chat_messenger/internal/protocol"
)

// UDPServer はUDPサーバーを表します。
type UDPServer struct {
	conn        *net.UDPConn
	roomManager chat.RoomManager
	userManager auth.UserManager
	port        string
}

// NewUDPServer は新しいUDPServerを生成します。
func NewUDPServer(port string, roomManager chat.RoomManager, userManager auth.UserManager) (*UDPServer, error) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("UDPアドレスの解決に失敗しました: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("UDPサーバーの起動に失敗しました: %w", err)
	}

	return &UDPServer{
		conn:        conn,
		roomManager: roomManager,
		userManager: userManager,
		port:        port,
	}, nil
}

// Start はUDPサーバーを起動し、クライアントからのメッセージを待ち受けます。
func (s *UDPServer) Start() error {
	fmt.Printf("UDPサーバーを起動しました (ポート: %s)...\n", s.port)

	// 現時点では受信処理は実装せず、サーバーの起動のみ行う
	// 実際の処理は後で実装する

	s.handleConnection(s.conn)
	return nil
}

// Close はUDPサーバーを停止します。
func (s *UDPServer) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *UDPServer) handleConnection(conn *net.UDPConn) {
	defer conn.Close()
	for {
		buf := make([]byte, 4096)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("Received from %v\n", remoteAddr)

		// バイトデータをUDPMessage構造体にデコード
		udpMessage, err := protocol.DecodeUDPMessage(buf[:n])
		if err != nil {
			fmt.Println("クライアントメッセージをデコードできませんでした:", err)
			continue
		}

		// ヘッダーからサイズ情報を取得
		roomNameSize := udpMessage.Header.RoomNameSize
		tokenSize := udpMessage.Header.TokenSize

		// ボディからルーム名を取得
		roomName := string(udpMessage.Body[:roomNameSize])

		// ボディからトークンを取得
		tokenStart := roomNameSize
		tokenEnd := tokenStart + tokenSize
		token := string(udpMessage.Body[tokenStart:tokenEnd])

		// ボディからメッセージを取得
		message := string(udpMessage.Body[tokenEnd:])

		fmt.Printf("ルーム名: %s, トークン: %s, メッセージ: %s\n", roomName, token, message)

		// トークンの検証処理
		user, err := s.validateToken(token, roomName)
		if err != nil {
			fmt.Printf("トークン検証エラー: %v\n", err)
			// エラーメッセージをクライアントに返す処理を追加することも可能
			continue
		}

		fmt.Printf("ユーザー '%s' からのメッセージを受信しました\n", user.GetName())

		// ユーザーのUDPアドレスを更新
		if udpUser, ok := user.(*chat.SimpleUser); ok {
			udpUser.SetUDPAddr(remoteAddr)
		}

		// ルームの検索
		room, err := s.roomManager.FindRoom(roomName)
		if err != nil {
			fmt.Printf("ルーム '%s' が見つかりません\n", roomName)
			continue
		}

		// ルーム内の全ユーザーにメッセージをブロードキャスト
		s.broadcastToRoom(conn, room, message, user)

		// // クライアントに確認応答を返す
		// _, err = conn.WriteToUDP([]byte(message), remoteAddr)
		// if err != nil {
		// 	fmt.Println("Error writing to UDP:", err)
		// }
	}
}

// broadcastToRoom はルーム内の全ユーザーにメッセージをブロードキャストします
func (s *UDPServer) broadcastToRoom(conn *net.UDPConn, room chat.Room, message string, sender chat.User) {
	senderName := sender.GetName()
	senderToken := sender.GetToken()

	formattedMessage := fmt.Sprintf("%s> %s", senderName, message)
	messageBytes := []byte(formattedMessage)

	for _, user := range room.GetUsers() {
		// 送信者自身には送信しない
		if user.GetToken() == senderToken {
			continue
		}

		// ユーザーのUDPアドレスを取得
		udpUser, ok := user.(*chat.SimpleUser)
		if !ok || udpUser.GetUDPAddr() == nil {
			// UDPアドレスが設定されていないユーザーはスキップ
			continue
		}

		// メッセージを送信
		_, err := conn.WriteToUDP(messageBytes, udpUser.GetUDPAddr())
		if err != nil {
			fmt.Printf("ユーザー '%s' へのメッセージ送信に失敗: %v\n", user.GetName(), err)
		} else {
			fmt.Printf("ユーザー '%s' にメッセージを送信しました\n", user.GetName())
		}
	}
}

// トークンの検証
func (s *UDPServer) validateToken(token string, roomName string) (chat.User, error) {
	// トークンからユーザーを検索
	user, err := s.userManager.FindUser(token)
	if err != nil {
		return nil, fmt.Errorf("無効なトークン: %w", err)
	}

	// ルームを検索
	room, err := s.roomManager.FindRoom(roomName)
	if err != nil {
		return nil, fmt.Errorf("ルームが存在しません: %w", err)
	}

	// ユーザーがルームに所属しているか確認
	roomUsers := room.GetUsers()
	userFound := false

	for _, u := range roomUsers {
		if u.GetToken() == token {
			userFound = true
			break
		}
	}

	if !userFound {
		return nil, fmt.Errorf("ユーザーはこのルームに所属していません")
	}

	return user, nil
}
