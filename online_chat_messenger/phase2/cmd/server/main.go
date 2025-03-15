package main

import (
	"fmt"
	"os"

	"online_chat_messenger/internal/auth"
	"online_chat_messenger/internal/chat"
	"online_chat_messenger/internal/network"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用法: server <ポート番号>")
		os.Exit(1)
	}
	port := os.Args[1]

	roomManager := chat.NewSimpleRoomManager()
	userManager := auth.NewSimpleUserManager()

	server, err := network.NewTCPServer(port, roomManager, userManager)
	if err != nil {
		fmt.Printf("サーバーの起動に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer server.Close()

	if err := server.Start(); err != nil {
		fmt.Printf("サーバーの実行中にエラーが発生しました: %v\n", err)
		os.Exit(1)
	}
}
