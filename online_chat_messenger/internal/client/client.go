package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Start() {
	fmt.Println("client start")

	conn, err := net.Dial("udp", ":8088")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// ユーザーからの入力を受け付ける
	fmt.Println("メッセージを入力してください。")
	scanner := bufio.NewScanner(os.Stdin) // 標準入力を受け付けるスキャナ
	scanner.Scan()                        // １行分の入力を取得する. この処理で入力待ちに

	conn.Write([]byte(scanner.Text())) // Write()は[]byte
	fmt.Println("クライアント:メッセージ送信")
}
