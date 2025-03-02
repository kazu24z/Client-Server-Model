package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// バッファプールを定義（ポインタ版）
var bufferPool = sync.Pool{
	// Newフィールドは、プールが空のとき、Get()が呼ばれたときに紐づく関数を実行する
	New: func() interface{} { // どんな型でもプールできるようにinterface{}を使用
		buffer := make([]byte, 4096)
		return &buffer // ポインタを返す
	},
}

func Start() {
	fmt.Println("server Start")

	udpAddr, err := net.ResolveUDPAddr("udp", ":8088")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// クライアントを管理するインスタンス
	clientManager := NewClientManager(300 * time.Second)

	// 期限切れの接続をマネージャから削除する処理を実行する
	go func() {
		ticker := time.NewTicker(60 * time.Second) // Ticker構造体生成時に、チャネルにデータを入れる
		defer ticker.Stop()

		for range ticker.C {
			// タイムアウト接続を除去する処理
			clientManager.RemoveInactiveClients()
		}
	}()

	handleConn(conn, clientManager)
}

func handleConn(conn *net.UDPConn, clientManager *ClientManager) {
	buf := make([]byte, 4096)

	for {
		// clientAddrがあるのは、UDPは接続状態を持たないため。どこから来たかをここで保持しておく
		n, clientAddr, err := conn.ReadFromUDP(buf) // こいつ、ブロッキングしてるからbufにデータが入るまでここで待機する。だから呼び出したとき、処理が終わらない

		if err != nil {
			log.Println("読み取りエラー:", err)
			continue // エラーがあったとしても、ループ自体は終わらせない = サーバは待機したまま
		}

		// プールからバッファポインタを取得
		bufPtr := bufferPool.Get().(*[]byte) // bufferPool.Get()はinterface{}型
		// 必要なサイズに調整してコピー
		bufCopy := (*bufPtr)[:n]
		copy(bufCopy, buf[:n]) // プールさせた[]byteスライスの実態にUDPで読み取ったデータをコピー

		// addrがあれば、サーバ側でクライアントを覚えておく
		clientManager.AddNewClient(clientAddr)

		// goroutine呼び出しでmainスレッドはメッセージ取得に専念できる
		go func(data []byte, bufPtr *[]byte, clientAddr *net.UDPAddr) {
			// 処理が終わったらポインタを返却
			defer bufferPool.Put(bufPtr)

			userName, messageBody := parseMessage(data, len(data))

			// userName, messageBodyが空のとき、データを捨てる
			if userName == "" || messageBody == "" {
				return
			}

			message := userName + ": " + messageBody

			// メッセージ送信
			clientManager.BroadCast(message, clientAddr, conn)

			fmt.Printf("[%s]: %s\n", clientAddr, message)
		}(bufCopy, bufPtr, clientAddr)
	}
}

func parseMessage(buf []byte, n int) (string, string) {

	if n < 1 {
		return "", ""
	}

	userNameLen := int(buf[0])

	if n < userNameLen+1 {
		return "", ""
	}

	userName := string(buf[1 : userNameLen+1])
	message := string(buf[userNameLen+1 : n])

	return userName, message

}
