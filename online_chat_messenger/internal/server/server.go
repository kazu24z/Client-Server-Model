package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

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

	handleConn(conn)
}

func handleConn(conn *net.UDPConn) {
	buf := make([]byte, 4096)

	// クライアントを管理するインスタンス
	clientManager := NewClientManager(300 * time.Second)

	for {
		// addrがあるのは、UDPは接続状態を持たないため。どこから来たかをここで保持しておく
		n, addr, err := conn.ReadFromUDP(buf) // こいつ、ブロッキングしてるからbufにデータが入るまでここで待機する。だから呼び出したとき、処理が終わらない

		if err != nil {
			log.Println("読み取りエラー:", err)
			continue // エラーがあったとしても、ループ自体は終わらせない = サーバは待機したまま
		}

		// addrがあれば、サーバ側でクライアントを覚えておく
		clientManager.AddNewClient(addr)

		// buf[:n] これで送信メッセージのbyte列が取得できる
		// buf[0] これはユーザー名を表現するバイト数が入っている → intにする
		// message = buf[:n] - buf[0] - buf[1:n+1]
		userName, messageBody := parseMessage(buf, n)

		// userName, messageBodyが空のとき、データを捨てる
		if userName == "" || messageBody == "" {
			continue
		}

		message := userName + ": " + messageBody

		// メッセージ送信
		clientManager.BroadCast(message, addr, conn)

		fmt.Printf("[%s]: %s\n", addr, message)
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
