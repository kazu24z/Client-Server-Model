package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// TCRPHeader はTCRPヘッダーを表します。
type TCRPHeader struct {
	RoomNameSize         uint8
	Operation            uint8
	State                uint8
	OperationPayloadSize uint8
}

// TCRPMessage はTCRPメッセージを表します。
type TCRPMessage struct {
	Header TCRPHeader
	Body   []byte
}

// EncodeTCRPMessage はTCRPメッセージをバイト列にエンコードします。
func EncodeTCRPMessage(msg TCRPMessage) ([]byte, error) {
	buf := new(bytes.Buffer) // new()は型を受け取り、その型でメモリ割り当てをする。値はゼロ値を適用

	// ヘッダーを書き込み
	if err := binary.Write(buf, binary.LittleEndian, msg.Header.RoomNameSize); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.Header.Operation); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.Header.State); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.Header.OperationPayloadSize); err != nil {
		return nil, err
	}

	// ボディを書き込み

	if _, err := buf.Write(msg.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeTCRPMessage はバイト列をTCRPメッセージにデコードします。
func DecodeTCRPMessage(data []byte) (TCRPMessage, error) {
	var msg TCRPMessage
	buf := bytes.NewReader(data)

	// ヘッダーを読み込み
	// binary.Read 関数は、io.Reader の現在の読み込み位置からデータを読み込み、読み込み後、その位置を進める
	if err := binary.Read(buf, binary.LittleEndian, &msg.Header.RoomNameSize); err != nil {
		return msg, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Header.Operation); err != nil {
		return msg, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Header.State); err != nil {
		return msg, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Header.OperationPayloadSize); err != nil {
		return msg, err
	}

	// ボディを読み込み
	bodySize := int(msg.Header.OperationPayloadSize)
	if bodySize > 0 {
		msg.Body = make([]byte, bodySize)
		if _, err := io.ReadFull(buf, msg.Body); err != nil {
			return msg, fmt.Errorf("ボディの読み込みエラー: %w", err)
		}
	} else {
		// ヘッダーサイズを計算
		headerSize := binary.Size(msg.Header)
		// 残りのデータをすべて読み込む
		remainingSize := len(data) - headerSize
		if remainingSize > 0 {
			msg.Body = make([]byte, remainingSize)
			if _, err := io.ReadFull(buf, msg.Body); err != nil {
				return msg, fmt.Errorf("残りデータの読み込みエラー: %w", err)
			}
		}
	}

	// デバッグ情報を出力
	fmt.Printf("デコード結果: ヘッダー=%+v, ボディ長=%d\n", msg.Header, len(msg.Body))
	if len(msg.Body) > 0 {
		fmt.Printf("ボディ内容: %s\n", string(msg.Body))
	}

	return msg, nil
}
