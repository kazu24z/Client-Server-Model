package network

import (
	"encoding/json"
	"fmt"
	"net"

	"online_chat_messenger/internal/auth"
	"online_chat_messenger/internal/chat"
	"online_chat_messenger/internal/protocol"
)

// ClientRequest はクライアントからのリクエストを表します。
type ClientRequest struct {
	RoomName  string `json:"room_name"`
	Password  string `json:"password,omitempty"`
	UserName  string `json:"user_name"`
	Operation uint8  // protocol/tcrp.go の operationと対応させる
	State     uint8  // protocol/tcrp.go の stateと対応させる
}

// TCPServer はTCPサーバーを表します。
type TCPServer struct {
	listener    net.Listener
	roomManager chat.RoomManager
	userManager auth.UserManager
}

// NewTCPServer は新しいTCPServerを生成します。
func NewTCPServer(port string, roomManager *chat.SimpleRoomManager, userManager *auth.SimpleUserManager) (*TCPServer, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("TCPサーバーの起動に失敗しました: %w", err)
	}
	return &TCPServer{
		listener:    listener,
		roomManager: roomManager,
		userManager: userManager,
	}, nil
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

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	fmt.Println("ユーザーリクエストバイト数（State0）", n)
	if err != nil {
		fmt.Printf("データの受信に失敗しました: %v\n", err)
		return
	}

	// TCRPメッセージをデコードする
	tcrpMsg, err := protocol.DecodeTCRPMessage(buffer[:n])
	if err != nil {
		fmt.Printf("TCRPメッセージのデコードに失敗しました: %v\n", err)
		return
	}

	// ペイロードをJSONにデコードする
	var request ClientRequest
	err = json.Unmarshal(tcrpMsg.Body, &request)
	if err != nil {
		fmt.Printf("JSONのデコードに失敗しました.State0（リクエスト）: %v\n", err)
		return
	}
	request.Operation = tcrpMsg.Header.Operation
	request.State = tcrpMsg.Header.State

	// リクエストの種類に応じて処理を分岐する
	switch {
	case request.Operation == 1 && request.State == 0: // チャットルーム作成リクエスト (初期化)
		s.handleCreateRoomRequest(conn, request)
	case request.Operation == 2 && request.State == 0: // チャットルーム参加リクエスト (初期化)
		s.handleJoinRoomRequest(conn, request)
	default:
		fmt.Println("不明なリクエストです")
	}
}

// handleCreateRoomRequest はクライアントからのルーム作成リクエストを処理します。
func (s *TCPServer) handleCreateRoomRequest(conn net.Conn, request ClientRequest) {
	fmt.Printf("ルーム作成リクエストを受けました: %+v\n", request)

	// ステータスコードを含むJSONペイロードを生成
	statusPayload, err := json.Marshal(map[string]int{"status": 0}) // 成功
	if err != nil {
		fmt.Printf("JSONのエンコードに失敗しました: %v\n", err)
		return
	}

	// リクエストの応答 (1)
	responseTCRPMessage := protocol.TCRPMessage{
		Header: protocol.TCRPHeader{
			Operation: 1,
			State:     1,
		},
		Body: statusPayload,
	}
	encodedResponseMessage, err := protocol.EncodeTCRPMessage(responseTCRPMessage)
	if err != nil {
		fmt.Printf("TCRPメッセージのエンコードに失敗しました: %v\n", err)
		return
	}

	fmt.Println("state1のエンコード済みデータ:", encodedResponseMessage)
	_, err = conn.Write(encodedResponseMessage)
	if err != nil {
		fmt.Printf("データの送信に失敗しました: %v\n", err)
		return
	}

	// トークンを生成
	token := auth.GenerateToken()

	// チャットルームを作成し、ホストを設定
	room, err := s.roomManager.CreateRoom(request.RoomName, request.Password)
	if err != nil {
		// ルーム作成失敗時の処理
		return
	}

	user := chat.NewUser(request.UserName, token, conn.RemoteAddr().String())

	err = room.AddUser(user, true) //trueでhostとして設定
	if err != nil {
		// ユーザー追加失敗時の処理
		return
	}
	s.userManager.RegisterUser(token, user)

	// トークンを含むJSONペイロードを生成
	tokenPayload, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		fmt.Printf("JSONのエンコードに失敗しました: %v\n", err)
		return
	}

	// リクエストの完了 (2)
	completeTCRPMessage := protocol.TCRPMessage{
		Header: protocol.TCRPHeader{
			Operation: 1,
			State:     2,
		},
		Body: tokenPayload,
	}
	encodedCompleteMessage, err := protocol.EncodeTCRPMessage(completeTCRPMessage)
	if err != nil {
		fmt.Printf("TCRPメッセージのエンコードに失敗しました: %v\n", err)
		return
	}

	_, err = conn.Write(encodedCompleteMessage)
	if err != nil {
		fmt.Printf("データの送信に失敗しました: %v\n", err)
		return
	}
}

// handleJoinRoomRequest はクライアントからのルーム参加リクエストを処理します。
func (s *TCPServer) handleJoinRoomRequest(conn net.Conn, request ClientRequest) {
	fmt.Printf("ルーム参加リクエストを受けました: %+v\n", request)

	// ステータスコードを含むJSONペイロードを生成
	statusPayload, err := json.Marshal(map[string]int{"status": 0}) // 成功
	if err != nil {
		fmt.Printf("JSONのエンコードに失敗しました: %v\n", err)
		return
	}

	// リクエストの応答 (1)
	responseTCRPMessage := protocol.TCRPMessage{
		Header: protocol.TCRPHeader{
			Operation: 2,
			State:     1,
		},
		Body: statusPayload,
	}
	encodedResponseMessage, err := protocol.EncodeTCRPMessage(responseTCRPMessage)
	if err != nil {
		fmt.Printf("TCRPメッセージのエンコードに失敗しました: %v\n", err)
		return
	}

	fmt.Println("state1のエンコード済みデータ:", encodedResponseMessage)
	_, err = conn.Write(encodedResponseMessage)
	if err != nil {
		fmt.Printf("データの送信に失敗しました: %v\n", err)
		return
	}

	// チャットルームを検索
	room, err := s.roomManager.FindRoom(request.RoomName)
	if err != nil {
		// ルームが見つからない場合のエラー処理
		fmt.Printf("ルームが見つかりませんでした: %v\n", err)
		return
	}

	// パスワードが一致するか確認（必要な場合）
	// ...

	// トークンを生成
	token := auth.GenerateToken()

	// ユーザーを作成
	user := chat.NewUser(request.UserName, token, conn.RemoteAddr().String())

	// チャットルームに参加
	err = room.AddUser(user, false) //falseでhostではない
	if err != nil {
		// 参加失敗時のエラー処理
		fmt.Printf("ルームへの参加に失敗しました: %v\n", err)
		return
	}

	// ユーザーを登録
	s.userManager.RegisterUser(token, user)

	// トークンを含むJSONペイロードを生成
	tokenPayload, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		fmt.Printf("JSONのエンコードに失敗しました: %v\n", err)
		return
	}

	// リクエストの完了 (2)
	completeTCRPMessage := protocol.TCRPMessage{
		Header: protocol.TCRPHeader{
			Operation: 2,
			State:     2,
		},
		Body: tokenPayload,
	}
	encodedCompleteMessage, err := protocol.EncodeTCRPMessage(completeTCRPMessage)
	if err != nil {
		fmt.Printf("TCRPメッセージのエンコードに失敗しました: %v\n", err)
		return
	}

	_, err = conn.Write(encodedCompleteMessage)
	if err != nil {
		fmt.Printf("データの送信に失敗しました: %v\n", err)
		return
	}
}

// Close はTCPサーバーを停止します。
func (s *TCPServer) Close() error {
	return s.listener.Close()
}
