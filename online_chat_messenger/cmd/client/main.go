package main

import (
	"log"
	"online_chat_messenger/internal/client"
)

func main() {
	c := client.NewClient()
	defer c.Close()

	if err := c.SetupTerminal(); err != nil {
		log.Fatalf("ターミナルのセットアップに失敗しました: %v", err)
	}

	c.RegisterUserName()

	if err := c.ConnectToServer("udp", ":8088"); err != nil {
		log.Fatalf("サーバーへの接続に失敗しました: %v", err)
	}

	c.AddMessage("チャットに接続しました。メッセージを入力してください。")

	go c.ReceiveMessages()

	c.HandleInput()
}
