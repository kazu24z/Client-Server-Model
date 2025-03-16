package protocol

import (
	"fmt"
)

type UDPHeader struct {
	RoomNameSize uint8
	TokenSize    uint8
}

type UDPMessage struct {
	Header UDPHeader
	Body   []byte
}

func EncodeUDPMessage(msg UDPMessage) ([]byte, error) {
	var encoded []byte

	return encoded, nil
}

func DecodeUDPMessage(data []byte) (UDPMessage, error) {
	if len(data) < 2 {
		return UDPMessage{}, fmt.Errorf("データが短すぎます")
	}

	header := UDPHeader{
		RoomNameSize: data[0],
		TokenSize:    data[1],
	}

	// ヘッダーの後のデータをボディとして扱う
	body := data[2:]

	// ボディの長さが期待される長さと一致するか確認
	expectedBodyLength := int(header.RoomNameSize) + int(header.TokenSize)
	if len(body) < expectedBodyLength {
		return UDPMessage{}, fmt.Errorf("ボディデータが不足しています")
	}

	return UDPMessage{
		Header: header,
		Body:   body,
	}, nil
}
