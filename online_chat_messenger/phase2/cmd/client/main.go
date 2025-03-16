package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"online_chat_messenger/internal/protocol"
)

// ユーザー入力を取得する関数
func getUserInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1] // 改行を削除
}

// サーバーに接続する関数
func connectToServer() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		return nil, fmt.Errorf("サーバーへの接続に失敗しました: %v", err)
	}
	return conn, nil
}

// ルーム作成リクエストを送信する関数
func sendRequest(conn net.Conn, choice, roomName, userName string) error {
	request := map[string]string{
		"room_name": roomName,
		"user_name": userName,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("JSONのエンコードに失敗しました: %v", err)
	}

	operation, err := strconv.ParseUint(choice, 10, 8)
	if err != nil {
		return fmt.Errorf("選択の変換に失敗しました: %v", err)
	}

	requestTCRPMessage := protocol.TCRPMessage{
		Header: protocol.TCRPHeader{
			Operation: uint8(operation),
			State:     0,
		},
		Body: requestBody,
	}

	encodedRequest, err := protocol.EncodeTCRPMessage(requestTCRPMessage)
	if err != nil {
		return fmt.Errorf("TCRPメッセージのエンコードに失敗しました: %v", err)
	}

	_, err = conn.Write(encodedRequest)
	if err != nil {
		return fmt.Errorf("サーバーへの送信に失敗しました: %v", err)
	}

	return nil
}

// サーバーからの応答を受信する関数
func receiveResponse(conn net.Conn) (protocol.TCRPMessage, error) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return protocol.TCRPMessage{}, fmt.Errorf("サーバーからの受信に失敗しました: %v", err)
	}

	responseTCRPMessage, err := protocol.DecodeTCRPMessage(buffer[:n])
	if err != nil {
		return protocol.TCRPMessage{}, fmt.Errorf("TCRPメッセージのデコードに失敗しました: %v", err)
	}

	return responseTCRPMessage, nil
}

func main() {
	// サーバーに接続
	conn, err := connectToServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// ユーザー入力を取得
	choice := getUserInput(reader, "選択してください（1: 新規ルーム作成, 2: 既存ルーム入室）: ")
	roomName := getUserInput(reader, "ルーム名を入力してください: ")
	userName := getUserInput(reader, "ユーザー名を入力してください: ")

	// ルーム作成/参加リクエストを送信
	err = sendRequest(conn, choice, roomName, userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 準拠応答を受信 (State = 1)
	responseTCRPMessage, err := receiveResponse(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if responseTCRPMessage.Header.State != 1 {
		fmt.Println("リクエスト処理に失敗しました。状態コード:", responseTCRPMessage.Header.State)
		return
	}

	fmt.Println("サーバーからの応答を受信しました。状態:", responseTCRPMessage.Header.State)

	// 完了応答を受信 (State = 2)
	completeTCRPMessage, err := receiveResponse(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if completeTCRPMessage.Header.State != 2 {
		fmt.Println("リクエスト完了に失敗しました。状態コード:", completeTCRPMessage.Header.State)
		return
	}

	fmt.Println("サーバーからの完了応答を受信しました。状態:", completeTCRPMessage.Header.State)

	// 応答を表示
	var response map[string]string
	err = json.Unmarshal(completeTCRPMessage.Body, &response)
	if err != nil {
		fmt.Printf("JSONのデコードに失敗しました: %v\n", err)
		return
	}

	if completeTCRPMessage.Header.Operation == 1 {
		fmt.Println("ルーム作成に成功しました！")
	} else if completeTCRPMessage.Header.Operation == 2 {
		fmt.Println("ルームへの参加に成功しました！")
	}

	token := response["token"] // サーバから返されたトークン
	fmt.Println("トークン:", token)

	// roomNameがレスポンスに含まれている場合は更新
	if respRoomName, ok := response["roomName"]; ok {
		roomName = respRoomName
	}
	fmt.Println("ルーム名:", roomName)

	// TCPの接続を閉じる
	conn.Close()

	// ===== UDPのチャットルーム処理に移行 =============
	udpConn, err := connectToServerUDP()
	if err != nil {
		fmt.Println("UDP接続に失敗しました。", err)
		return
	}
	defer udpConn.Close()

	// 受信処理をゴルーチンで実行
	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := udpConn.Read(buf)
			if err != nil {
				fmt.Println("サーバからの受信に失敗しました:", err)
				return
			}
			formatReceiveMessage(buf[:n])
		}
	}()

	// メインスレッドで送信処理を実行
	reader = bufio.NewReader(os.Stdin)
	for {
		message := getUserInput(reader, userName+"> ")
		if message == "/exit" {
			fmt.Println("チャットを終了します")
			break
		}
		sendChat(udpConn, message, token, roomName)
	}
}

// UDPでサーバーに接続する関数
func connectToServerUDP() (*net.UDPConn, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", ":8089")
	if err != nil {
		return nil, fmt.Errorf("UDPアドレスの解決に失敗しました: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, fmt.Errorf("UDPサーバーへの接続に失敗しました: %v", err)
	}

	return conn, nil
}

func sendChat(conn net.Conn, message string, token string, roomName string) {
	// UDPメッセージのプロトコルに則ってデータを用意する
	roomNameBytes := []byte(roomName)
	tokenBytes := []byte(token)
	messageBytes := []byte(message)

	udpHeader := protocol.UDPHeader{
		RoomNameSize: uint8(len(roomNameBytes)),
		TokenSize:    uint8(len(tokenBytes)),
	}

	// ボディ部分を構築: roomName + token + message の順
	bodySize := len(roomNameBytes) + len(tokenBytes) + len(messageBytes)
	body := make([]byte, bodySize)

	// ルーム名をコピー
	offset := 0
	copy(body[offset:], roomNameBytes)
	offset += len(roomNameBytes)

	// トークンをコピー
	copy(body[offset:], tokenBytes)
	offset += len(tokenBytes)

	// メッセージをコピー
	copy(body[offset:], messageBytes)

	// 送信データを構築（ヘッダー + ボディ）
	data := make([]byte, 2+len(body))
	data[0] = udpHeader.RoomNameSize
	data[1] = udpHeader.TokenSize
	copy(data[2:], body)

	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("UDPサーバへの送信に失敗しました:", err)
		return
	}
}

func formatReceiveMessage(buf []byte) {
	message := string(buf)

	// 画面をクリアせずに、現在の入力行を消去して新しいメッセージを表示
	fmt.Print("\r\033[K") // カーソルを行頭に移動して行をクリア
	fmt.Println(message)  // メッセージを表示

	// 入力プロンプトを再表示
	fmt.Print("\r")
}
