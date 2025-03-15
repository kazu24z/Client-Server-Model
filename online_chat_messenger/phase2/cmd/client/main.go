package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	fmt.Println("受信バイト数:", n)
	if err != nil && !errors.Is(err, io.EOF) {
		return protocol.TCRPMessage{}, fmt.Errorf("サーバーからの受信に失敗しました: %v", err)
	}

	responseTCRPMessage, err := protocol.DecodeTCRPMessage(buffer[:n])
	if err != nil && !errors.Is(err, io.EOF) {
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

	// ルーム作成リクエストを送信
	err = sendRequest(conn, choice, roomName, userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 最初の応答を受信
	responseTCRPMessage, err := receiveResponse(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if responseTCRPMessage.Header.State != 1 {
		fmt.Println("ルーム作成に失敗しました:", responseTCRPMessage.Header.State)
		return
	}

	// 完了応答を受信
	completeTCRPMessage, err := receiveResponse(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if completeTCRPMessage.Header.State != 2 {
		fmt.Println("ルーム作成に失敗しました")
		return
	}

	// 応答を表示
	var response map[string]string
	err = json.Unmarshal(completeTCRPMessage.Body, &response)
	if err != nil {
		fmt.Printf("JSONのデコードに失敗しました: %v\n", err)
		return
	}
	if responseTCRPMessage.Header.Operation == 1 {
		fmt.Printf("ルーム作成に成功しました！")
	} else if responseTCRPMessage.Header.Operation == 2 {
		fmt.Printf("ルームの参加に成功しました")
	}
}
