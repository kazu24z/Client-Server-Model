package server

import (
	"fmt"
	"log"
	"net"
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
	for {
		n, addr, err := conn.ReadFromUDP(buf) // こいつ、ブロッキングしてるからbufにデータが入るまでここで待機する。だから呼び出したとき、処理が終わらない
		if err != nil {
			log.Println("読み取りエラー:", err)
			continue // エラーがあったとしても、ループ自体は終わらせない = サーバは待機したまま
		}
		fmt.Printf("メッセージ受信 [%s]: %s\n", addr, string(buf[:n]))
	}
}
