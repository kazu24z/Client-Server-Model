package client

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"

	"golang.org/x/term"
)

type Client struct {
	conn       net.Conn
	userName   string
	inputLine  string
	messageLog []string
	oldState   *term.State
	mutex      sync.Mutex
}

func NewClient() *Client {
	return &Client{
		messageLog: make([]string, 0),
	}
}

func (c *Client) SetupTerminal() error {
	// ターミナルをRawモードに設定
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("ターミナルの設定に失敗しました: %w", err)
	}

	// 画面をクリア
	cmd := exec.Command("clear") // os/execパッケージ...シェルコマンドをGoで実行できる
	cmd.Stdout = os.Stdout       // Goに渡したコマンドの出力を、os.Stdoutつまり標準出力に指定
	err = cmd.Run()              // "clearコマンド実行"
	if err != nil {
		return fmt.Errorf("画面のクリアに失敗しました: %w", err)
	}

	c.oldState = oldState
	return nil
}

func (c *Client) RegisterUserName() {
	fmt.Print("ユーザー名: ")
	input := ""
	for {
		char := make([]byte, 1)
		os.Stdin.Read(char)

		if char[0] == 13 { // Enter key
			if input != "" {
				break
			}
			fmt.Print("\r\nユーザー名が空です。再入力してください: ")
			continue
		} else if char[0] == 127 || char[0] == 8 { // Backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b") // バックスペース、スペース、バックスペース
			}
		} else {
			input += string(char)
			fmt.Print(string(char))
		}
	}
	c.userName = input
	fmt.Print("\r\n")
	c.renderScreen()
}

func (c *Client) ConnectToServer(network string, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return fmt.Errorf("サーバーへの接続に失敗しました: %w", err)
	}

	c.conn = conn
	return nil
}

func (c *Client) AddMessage(message string) {
	c.messageLog = append(c.messageLog, message)
	c.renderScreen()
}

func (c *Client) ReceiveMessages() {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			c.mutex.Lock()
			c.messageLog = append(c.messageLog, "エラー: サーバーとの接続が切断されました")
			c.renderScreen()
			c.mutex.Unlock()
			log.Printf("受信エラー: %v", err)
			return
		}

		message := string(buf[:n])

		c.mutex.Lock()
		// 受信メッセージをログに追加
		c.messageLog = append(c.messageLog, message)
		// 画面を再描画
		c.renderScreen()
		c.mutex.Unlock()
	}
}

func (c *Client) HandleInput() {
	for {
		char := make([]byte, 1)
		_, err := os.Stdin.Read(char)
		if err != nil {
			c.mutex.Lock()
			c.messageLog = append(c.messageLog, "エラー: 入力の読み取りに失敗しました")
			c.renderScreen()
			c.mutex.Unlock()
			log.Printf("入力エラー: %v", err)
			return
		}

		c.mutex.Lock()
		if char[0] == 13 { // Enter key
			// Enterキーが押されたらメッセージを送信
			if c.inputLine != "" {
				// 送信メッセージを作成
				message := make([]byte, 1)
				userNameLen := len([]byte(c.userName))
				message[0] = byte(userNameLen)
				message = append(message, []byte(c.userName)...)
				message = append(message, []byte(c.inputLine)...)

				// 自分のメッセージをログに追加
				c.messageLog = append(c.messageLog, c.userName+": "+c.inputLine)

				// メッセージ送信
				_, err := c.conn.Write(message)
				if err != nil {
					c.messageLog = append(c.messageLog, "エラー: メッセージの送信に失敗しました")
					log.Printf("送信エラー: %v", err)
				}

				// 入力行をクリア
				c.inputLine = ""
			}
		} else if char[0] == 127 || char[0] == 8 { // Backspace
			// 文字を1つ削除
			if len(c.inputLine) > 0 {
				c.inputLine = c.inputLine[:len(c.inputLine)-1]
			}
		} else if char[0] == 3 { // Ctrl+C
			term.Restore(int(os.Stdin.Fd()), c.oldState)
			os.Exit(0)
		} else {
			// 文字を追加
			c.inputLine += string(char)
		}

		c.renderScreen()
		c.mutex.Unlock()
	}
}

func (c *Client) renderScreen() {
	// 画面をクリアしてカーソルを左上に移動
	fmt.Print("\033[2J\033[H")

	// メッセージログを表示（左揃えで）
	for _, msg := range c.messageLog {
		// 必ず左揃えになるようにクリアしてから表示
		fmt.Print("\r")
		fmt.Println(msg)
	}

	// 最後に入力プロンプトと現在の入力内容を表示
	fmt.Print("\r" + c.userName + "> " + c.inputLine)
}

func (c *Client) Close() {
	if c.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), c.oldState)
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
