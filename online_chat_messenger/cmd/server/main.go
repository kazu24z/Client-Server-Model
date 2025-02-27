package main

import (
	"fmt"
	// ここにinternalで定義したパッケージをインポート
	"online_chat_messenger/internal/server"
)

func main() {
	fmt.Println("run server")

	server.Start()
}
