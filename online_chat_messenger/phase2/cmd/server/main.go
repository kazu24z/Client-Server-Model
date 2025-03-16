package main

import (
	"fmt"
	"os"

	"online_chat_messenger/internal/auth"
	"online_chat_messenger/internal/chat"
	"online_chat_messenger/internal/network"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("使用法: server <ポート番号>")
		os.Exit(1)
	}
	tcpPort := os.Args[1]
	udpPort := os.Args[2]

	roomManager := chat.NewSimpleRoomManager()
	userManager := auth.NewSimpleUserManager()

	// TCPサーバーの初期化
	tcpServer, err := network.NewTCPServer(tcpPort, roomManager, userManager)
	if err != nil {
		fmt.Printf("TCPサーバーの起動に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer tcpServer.Close()

	// UDPサーバーの初期化（同じポートを使用）
	udpServer, err := network.NewUDPServer(udpPort, roomManager, userManager)
	if err != nil {
		fmt.Printf("UDPサーバーの起動に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer udpServer.Close()

	// UDPサーバーを別のゴルーチンで起動
	go func() {
		if err := udpServer.Start(); err != nil {
			fmt.Printf("UDPサーバーの実行中にエラーが発生しました: %v\n", err)
			os.Exit(1)
		}
	}()

	// TCPサーバーをメインゴルーチンで起動
	if err := tcpServer.Start(); err != nil {
		fmt.Printf("TCPサーバーの実行中にエラーが発生しました: %v\n", err)
		os.Exit(1)
	}
}
